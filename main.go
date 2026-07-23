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
	LoggingConf := &LoggerConf{} 	

	// logger parameters config
	LogParamConf := &LoggerParamConf{
		Reqclient: client,
		Req: req,
		ReqBodySize: bodySize,
		AllowReqBody: allowedBody,
	}

	LoggingConf = Log(LogParamConf)

	if LoggingConf.Reqbody == "" {
		fmt.Println("body is empty")
	}

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer LoggingConf.Req.Body.Close() 	

	if *stream {
		fmt.Println("streaming live data:")
		_, err := io.Copy(os.Stdout, LoggingConf.Req.Body)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else {

		// checking if request failed if yes then log
		if LoggingConf.Req.StatusCode >= 400 { // anything over 400 means request wasnt successful
			fmt.Println("request failed with status:", LoggingConf.Req.Status)
		}

		// outputting

		fmt.Println("status:", LoggingConf.Req.Status)
		fmt.Println("latency:", LoggingConf.Reqtime) // printing latency

		fmt.Println("\nheaders:")
		for k, v := range LoggingConf.Req.Header {
			fmt.Println(k+":", v) // key:value output style for headers
		}

		fmt.Println("\nbody:")

		// var format bytes.Buffer // pretty printed body will be stored here before outputted

		fmt.Println(LoggingConf.Reqbody)

		fmt.Println("logged request body with size", *LoggingConf.Reqlen)
	}
}
