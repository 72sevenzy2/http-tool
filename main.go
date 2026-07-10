package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// for composable header inputs
type HeaderFlags []string

func main() {
	var headers HeaderFlags	

	// cli inputs
	flag.Var(&headers, "H", "Header (key:value)")
	stream := flag.Bool("stream", false, "live response") // for streaming live response
	method := flag.String("x", "GET", "http method")
	allowedBody := flag.Bool("b", false, "allow request body logging.")
	bodySize := flag.Int("s", 0, "body size")
	session := flag.Bool("session", false, "enable session mode.")

	data := flag.String("d", "", "request data")

	flag.Parse()

	// before validating flag.Args for other features, check is session is set first.
	if *session {
		scanner := bufio.NewScanner(os.Stdin)
		store := NewStore()

		StartSession(scanner, store)
	}

	if err := Validate(flag.Args()); err != nil {
		fmt.Println(errors.New(err.Error()))
		return
	}

	url := flag.Args()[0]

	uc := strings.ToUpper(*method) // normalise request method input as all caps for easier comparison
	var v string                   // holding actual method after validating if method is appropriate

	// simplified
	allowed := map[string]bool{
		"GET":  true,
		"POST": true,
	}

	if allowed[uc] {
		v = uc
	} else {
		fmt.Println("method not allowed")
		return
	}

	// body reader for post data
	var body io.Reader

	if *data != "" {
		body = strings.NewReader(*data)
	}

	req, err := http.NewRequest(v, url, body)

	// req, err := http.NewRequest("GET", url, nil) // initialising an http request
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// add heders
	if err := AddHeaders(req, headers); err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // timeout so that request doesnt last forever
	} // initialising the client

	// debugging
	fmt.Println("bodyAllowed:", *allowedBody)

	// logging
	end, resp, bodyprev, bodysize, err := Log(client, req, allowedBody, bodySize)

	if bodyprev == "" {
		fmt.Println("body is empty")
	}

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer resp.Body.Close() // important as not closing resp.Body would lead to performance issues + leaks, aswell as its apart of the ReadCloser interface so it has be closed.

	if *stream {
		fmt.Println("streaming live data:")
		_, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else {

		// checking if request failed if yes then log
		if resp.StatusCode >= 400 { // anything over 400 means request wasnt successful
			fmt.Println("request failed with status:", resp.Status)
		}

		// outputting

		fmt.Println("status:", resp.Status)
		fmt.Println("latency:", end) // printing latency

		fmt.Println("\nheaders:")
		for k, v := range resp.Header {
			fmt.Println(k+":", v) // key:value output style for headers
		}

		fmt.Println("\nbody:")

		// var format bytes.Buffer // pretty printed body will be stored here before outputted

		fmt.Println(bodyprev)

		fmt.Println("logged request body with size", *bodysize)
	}
}
