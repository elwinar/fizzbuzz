# FizzBuzz API

## Endpoints

### Compute a Fizz Buzz [`GET /`]

Compute a fizz buzz using the provided parameters.

| parameter | type    | default | description                             |
| --------- | ------- | ------- | --------------------------------------- |
| int1      | integer | 3       | first value to look for multiples of    |
| int2      | integer | 5       | second value to look for multiples of   |
| limit     | integer | 100     | number of integers to consider          |
| str1      | string  | fizz    | value to replace multiples of int1 with |
| str2      | string  | buzz    | value to replace multiples of int2 with |

#### Examples

- *Nominal case*

	```
	GET /?int1=3&int2=5&str1=bizz&str2=buzz&limit=20 HTTP/1.1
	```

	```
	HTTP/1.1 200 OK
	Content-Length: 122
	Content-Type: application/json

	{
	    "result": [
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
		"buzz"
	    ]
	}
	```

- *Invalid limit*

	```
	GET /?limit=wrong HTTP/1.1
	```

	```
	HTTP/1.1 400 Bad Request
	Content-Length: 88
	Content-Type: application/json

	{
	    "error": "parsing \"limit\" parameter: strconv.Atoi: parsing \"wrong\": invalid syntax"
	}
	```

### Retrieve most frequent request [`GET /statistics`]

Retrieve the most frequent bizzbuzz request parameters and the total times it
has been used. Default values aren't differenciated from the same explicit
parameters.

####

- *Nominal case*

	```
	GET /statistics HTTP/1.1
	```

	```
	HTTP/1.1 200 OK
	Content-Length: 81
	Content-Type: application/json

	{
	    "request": {
		"int1": 3,
		"int2": 5,
		"limit": 100,
		"str1": "fizz",
		"str2": "buzz"
	    },
	    "total": 3
	}
	```
