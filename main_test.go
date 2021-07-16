package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func getT(t *testing.T, endpoint string) *http.Response {
	res, err := http.Get(endpoint)
	if err != nil {
		t.Fatalf(`get: %s`, err.Error())
	}
	return res
}

func readAllT(t *testing.T, src io.Reader) []byte {
	raw, err := io.ReadAll(src)
	if err != nil {
		t.Fatalf(`readAll: %s`, err.Error())
	}
	return raw
}

func unmarshalT(t *testing.T, src []byte, dest interface{}) {
	err := json.Unmarshal(src, dest)
	if err != nil {
		t.Fatalf(`unmarshal: %s`, err.Error())
	}
}

func fizzBuzzT(t *testing.T, srv *httptest.Server, req Request) FizzBuzzResponse {
	var out FizzBuzzResponse
	res := getT(t, fmt.Sprintf(`%s%s?%s`, srv.URL, "/", req.URLEncode()))
	raw := readAllT(t, res.Body)
	unmarshalT(t, raw, &out)
	return out
}

func statisticsT(t *testing.T, srv *httptest.Server) StatisticsResponse {
	var out StatisticsResponse
	res := getT(t, fmt.Sprintf(`%s%s`, srv.URL, "/statistics"))
	raw := readAllT(t, res.Body)
	unmarshalT(t, raw, &out)
	return out
}

func TestService_FizzBuzz(t *testing.T) {
	service := NewService()
	server := httptest.NewServer(service.Mux)

	for n, c := range map[string]struct {
		input    Request
		expected []string
	}{
		"default": {
			input: Request{
				Int1:  3,
				Int2:  5,
				Str1:  "fizz",
				Str2:  "buzz",
				Limit: 100,
			},
			expected: []string{
				"1",
				"2",
				"fizz",
				"4",
				"buzz",
				"fizz",
				"7",
				"8",
				"fizz",
				"buzz",
				"11",
				"fizz",
				"13",
				"14",
				"fizzbuzz",
				"16",
				"17",
				"fizz",
				"19",
				"buzz",
				"fizz",
				"22",
				"23",
				"fizz",
				"buzz",
				"26",
				"fizz",
				"28",
				"29",
				"fizzbuzz",
				"31",
				"32",
				"fizz",
				"34",
				"buzz",
				"fizz",
				"37",
				"38",
				"fizz",
				"buzz",
				"41",
				"fizz",
				"43",
				"44",
				"fizzbuzz",
				"46",
				"47",
				"fizz",
				"49",
				"buzz",
				"fizz",
				"52",
				"53",
				"fizz",
				"buzz",
				"56",
				"fizz",
				"58",
				"59",
				"fizzbuzz",
				"61",
				"62",
				"fizz",
				"64",
				"buzz",
				"fizz",
				"67",
				"68",
				"fizz",
				"buzz",
				"71",
				"fizz",
				"73",
				"74",
				"fizzbuzz",
				"76",
				"77",
				"fizz",
				"79",
				"buzz",
				"fizz",
				"82",
				"83",
				"fizz",
				"buzz",
				"86",
				"fizz",
				"88",
				"89",
				"fizzbuzz",
				"91",
				"92",
				"fizz",
				"94",
				"buzz",
				"fizz",
				"97",
				"98",
				"fizz",
				"buzz",
			},
		},
	} {
		t.Run(n, func(t *testing.T) {
			output := fizzBuzzT(t, server, c.input)
			if !reflect.DeepEqual(c.expected, output.Result) {
				t.Errorf(`unexpected result: wanted %+v, got %+v`, c.expected, output)
				return
			}
		})
	}
}

func TestService_Statistics(t *testing.T) {
	service := NewService()
	server := httptest.NewServer(service.Mux)

	var (
		req1 = Request{Int1: 3, Int2: 5, Str1: "fizz", Str2: "buzz", Limit: 10}
		req2 = Request{Int1: 5, Int2: 3, Str1: "fizz", Str2: "buzz", Limit: 100}
		req3 = Request{Int1: 1, Int2: 10, Str1: "na", Str2: "batman!", Limit: 100}
		req4 = Request{Int1: 2, Int2: 6, Str1: "foo", Str2: "bar", Limit: 157}
	)

	_ = fizzBuzzT(t, server, req1)

	_ = fizzBuzzT(t, server, req2)
	_ = fizzBuzzT(t, server, req2)
	_ = fizzBuzzT(t, server, req2)

	_ = fizzBuzzT(t, server, req3)
	_ = fizzBuzzT(t, server, req3)

	_ = fizzBuzzT(t, server, req4)
	_ = fizzBuzzT(t, server, req4)
	_ = fizzBuzzT(t, server, req4)
	_ = fizzBuzzT(t, server, req4)

	output := statisticsT(t, server)
	if output.Request != req4 {
		t.Errorf(`unexpected most frequent request: wanted %+v, got %+v`, req4, output.Request)
		return
	}

	_ = fizzBuzzT(t, server, req2)
	_ = fizzBuzzT(t, server, req2)

	output = statisticsT(t, server)
	if output.Request != req2 {
		t.Errorf(`unexpected most frequent request: wanted %+v, got %+v`, req2, output.Request)
		return
	}
}
