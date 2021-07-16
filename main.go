package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	var bind string
	flag.StringVar(&bind, "bind", ":8080", "address to bind on")
	flag.Parse()

	service := NewService()

	server := &http.Server{
		Addr:    bind,
		Handler: service.Mux,
	}

	go func() {
		signals := make(chan os.Signal, 2)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Println("shuting server down:", err.Error())
			return
		}
	}()

	log.Println("starting")
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("listening:", err.Error())
		os.Exit(1)
		return
	}
	log.Println("stopping")
}

// Service for the various endpoints and data that needs to persist throughout
// requests.
type Service struct {
	requests   chan Request
	counters   map[Request]int
	max        Request
	maxCounter int
	lock       sync.RWMutex

	Mux *http.ServeMux
}

// NewService initialize a new Service struct and start the aggregation
// routine used by the endpoints.
func NewService() *Service {
	s := Service{
		requests: make(chan Request),
		counters: make(map[Request]int),
		Mux:      http.NewServeMux(),
	}
	s.Mux.HandleFunc("/", s.FizzBuzz)
	s.Mux.HandleFunc("/statistics", s.Statistics)

	go s.Aggregate()

	return &s
}

// Request for a fizz-buzz.
type Request struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
	Limit int    `json:"limit"`
}

func (r Request) URLEncode() string {
	return url.Values{
		"int1":  []string{strconv.Itoa(r.Int1)},
		"int2":  []string{strconv.Itoa(r.Int2)},
		"str1":  []string{r.Str1},
		"str2":  []string{r.Str2},
		"limit": []string{strconv.Itoa(r.Limit)},
	}.Encode()
}

// parseInt from a request, with default fallback if the value is empty (or not
// found), returning an eventual parsing error.
func parseInt(r *http.Request, key string, def int) (int, error) {
	raw := r.FormValue(key)
	if len(raw) == 0 {
		return def, nil
	}
	res, err := strconv.Atoi(raw)
	if err != nil {
		return res, fmt.Errorf("parsing %q parameter: %w", key, err)
	}
	return res, nil
}

// ParseRequest from the raw HTTP Request, and validates the parameters for
// consistency.
func ParseRequest(r *http.Request) (req Request, err error) {
	req.Int1, err = parseInt(r, "int1", 3)
	if err != nil {
		return req, err
	}

	req.Int2, err = parseInt(r, "int2", 5)
	if err != nil {
		return req, err
	}

	req.Str1 = r.FormValue("str1")
	if len(req.Str1) == 0 {
		req.Str1 = "fizz"
	}

	req.Str2 = r.FormValue("str2")
	if len(req.Str2) == 0 {
		req.Str2 = "buzz"
	}

	req.Limit, err = parseInt(r, "limit", 100)
	if err != nil {
		return req, err
	}
	if req.Limit < 0 {
		return req, fmt.Errorf(`invalid %q parameter: must be a positive integer`, "limit")
	}

	return req, nil
}

type FizzBuzzResponse struct {
	Result []string `json:"result"`
}

// write the status and JSON payload to the ResponseWriter, defaults to 500 and
// an error description if the payload can't be marshalled.
func write(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")

	raw, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": %q}`, err.Error())))
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(raw)
}

// FizzBuzz endpoint computes the list of string for a set of parameters.
func (s *Service) FizzBuzz(w http.ResponseWriter, r *http.Request) {
	req, err := ParseRequest(r)
	if err != nil {
		write(w, http.StatusBadRequest, struct {
			Err string `json:"error"`
		}{
			Err: err.Error(),
		})
		return
	}

	res := make([]string, 0, req.Limit)
	for i := 1; i <= req.Limit; i++ {
		var v string
		switch {
		case i%req.Int1 == 0 && i%req.Int2 == 0:
			v = fmt.Sprintf("%s%s", req.Str1, req.Str2)
		case i%req.Int1 == 0:
			v = req.Str1
		case i%req.Int2 == 0:
			v = req.Str2
		default:
			v = strconv.Itoa(i)
		}
		res = append(res, v)
	}

	write(w, http.StatusOK, FizzBuzzResponse{
		Result: res,
	})
	s.requests <- req
}

// Aggregate keeps track of the number calls done for each request, and of the
// most frequent request.
func (s *Service) Aggregate() {
	for req := range s.requests {
		s.counters[req] += 1

		// We don't need to lock the loop itself, as the aggregation
		// routine is the only one that access the map itself. This
		// makes the lock mostly free when the most frequent request is
		// somewhat stable.
		if s.counters[req] > s.maxCounter {
			s.lock.Lock()
			s.max = req
			s.maxCounter = s.counters[req]
			s.lock.Unlock()
		}
	}
}

type StatisticsResponse struct {
	Request Request `json:"request"`
	Total   int     `json:"total"`
}

// Statistics endpoint returns the most frequent request and the number of
// times its been called.
func (s *Service) Statistics(w http.ResponseWriter, r *http.Request) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	write(w, http.StatusOK, StatisticsResponse{
		Request: s.max,
		Total:   s.maxCounter,
	})
}
