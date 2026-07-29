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
	// flag for concurrent requests
	Reqs := flag.Int("c", 1, "number of requests.") // will be passed in Log()

	data := flag.String("d", "", "request data")

	flag.Parse() // parse flags

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

	// logger parameters config
	LogParamConf := &LoggerParamConf{
		Reqclient:    client,
		Req:          req,
		ReqBodySize:  bodySize,
		AllowReqBody: allowedBody,
		ReqNo:        Reqs,
	}

	LoggingConf := Log(LogParamConf)

	// todo: make one big nested loop for loggingConf.Results instead of redoing.
	for _, result := range LoggingConf.Results {
		if result.Reqbody == "" {
			fmt.Println("body is empty")
		}

		if result.Reqerr != nil { // LoggingConf.Reqerr is a []error.
			LoggingConf.Err()
			return
		}
		if result.ReqResponse != nil {
			continue
		}
		defer result.ReqResponse.Body.Close()
		if *stream {
			fmt.Println("streaming live data:")
			_, err := io.Copy(os.Stdout, result.ReqResponse.Body)
			if err != nil {
				log.Fatal(err.Error())
				return
			}
		} else {

			// checking if request failed if yes then log
			if result.ReqResponse.StatusCode >= 400 { // anything over 400 means request wasnt successful
				fmt.Println("request failed with status:")
				LoggingConf.DisplayReqStatus()
			}

			// outputting
			fmt.Println("request status:")
			LoggingConf.DisplayReqStatus()

			LoggingConf.DisplayReqLatency()

			fmt.Println("\nheaders:")
			for k, v := range result.ReqResponse.Header {
				fmt.Println(k+":", v) // key:value output style for headers
			}

			fmt.Println("\nbody:")

			// var format bytes.Buffer // pretty printed body will be stored here before outputted

			LoggingConf.DisplayReqBody()

			LoggingConf.DisplayReqBodyLen()

		}
	}
}
