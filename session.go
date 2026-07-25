package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// start a interactive session
func StartSession(b *bufio.Scanner, store *Data) {
	fmt.Println("session started.")
	for {
		fmt.Print(">")
		isOk := b.Scan() // take input
		if !isOk {
			fmt.Println(b.Err()) // possible scanner error case
		}

		input := b.Text()
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue // skip current iteration if no input
		}

		upperInput := strings.ToUpper(parts[0]) // "VAR", "GET", "DEL", "EXIT"

		switch upperInput {
		// declare variables (can be to store headers, urls, etc)
		case "VAR":
			ok := HandleVar(parts, store)
			if !ok {
				continue
			}
		case "GET":
			ok := HandleGet(parts, store)
			if !ok {
				continue
			}

		case "DEL":
			ok := HandleDel(parts, store)
			if !ok {
				continue
			}

			// actual api testing logic (GET only for now)
		case "TEST":
			val, ok := HandleTestValdiation(parts, store)
			if !ok {
				continue
			}

			//  utility variables
			UtilConf := InitTestUtilConf()
			// request dependent variables

			ReqConf := InitTestReqConf()
			// default request (GET) so cl does not stay nil if user were not to pass -X
			ReqConf.cl, ReqConf.clErr = http.NewRequest(http.MethodGet, val, nil)
			// clErr is handled below, after flags are parsed and appropriate requests are made.

			// parse header arguments
			for i := range parts {
				upc := strings.ToUpper(parts[i]) // normalize "-h" to all uppercase

				if upc == "-H" {
					headerSlice, ok := HandleHeaderAppending(parts, i, upc)
					if !ok {
						UtilConf.pass = false
						continue
					}

					ReqConf.reqHeaders = headerSlice // store headerSLice contents which were extracted.
				}

				// to check if argument is for req methods
				if upc == "-X" {
					newReqMethod, exists := HandleRequestMethodValidation(parts, i)
					if !exists {
						UtilConf.pass = false
						continue
					}

					if newReqMethod == "POST" {
						UtilConf.reqType = http.MethodPost

						uppercased1, exists := HandlePostValidation(parts, i)
						if !exists {
							UtilConf.pass = false
							continue
						}

						if uppercased1 == "-D" { // for normal json data

							cleanedData, ok := HandleJsonDataInput(parts, i)
							if !ok {
								UtilConf.pass = false
								continue
							}

							// assign jsonData to cleaned
							UtilConf.jsonData = cleanedData

						}

						if uppercased1 == "-F" { // form mode
							formParts, ok := HandleFormDataInput(parts, i)
							if !ok {
								UtilConf.pass = false
								continue
							}

							// appending the form values
							data := url.Values{}

							data.Set(formParts[0], formParts[1])
							enc := data.Encode() // encode

							ReqConf.formBody = strings.NewReader(enc)
						}

					} else {
						if newReqMethod == "GET" {
							UtilConf.reqType = http.MethodGet
						} else {
							fmt.Println("invalid method type.")
							continue
						}
					}

				}

				// small flag value to ensure headers dont overwrite eachother
				isJsonH := false
				isFormH := false

				client := http.Client{}
				// cl, err := http.NewRequest(reqType, val, nil) // new request
				if UtilConf.jsonData != "" && UtilConf.reqType == http.MethodPost {
					ReqConf.cl, ReqConf.clErr = http.NewRequest(UtilConf.reqType, val, strings.NewReader(UtilConf.jsonData))

					// set content type to json after creating request
					ReqConf.cl.Header.Add("Content-Type", "application/json")
					isJsonH = true
				} else {
					if UtilConf.jsonData == "" && UtilConf.reqType == http.MethodPost { // meaning that its form related request body data
						ReqConf.cl, ReqConf.clErr = http.NewRequest(UtilConf.reqType, val, ReqConf.formBody)

						// setting appropriate header afterwards
						ReqConf.cl.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						isFormH = true
					}
					if UtilConf.reqType == http.MethodGet { // standard get method
						ReqConf.cl, ReqConf.clErr = http.NewRequest(http.MethodGet, val, nil)
					}

				}

				if ReqConf.clErr != nil {
					fmt.Println(ReqConf.clErr)
					continue
				}

				if !UtilConf.pass {
					continue
				}

				if !isJsonH && !isFormH {
					if len(ReqConf.reqHeaders) == 2 {
						ReqConf.cl.Header.Add(ReqConf.reqHeaders[0], ReqConf.reqHeaders[1])
					}
				}

				resp, err2 := client.Do(ReqConf.cl) // send the request to the url provided
				if err2 != nil {
					fmt.Println(err2.Error())
					continue
				}

				// shadow auth header from resp
				clonedH := resp.Header.Clone()
				clonedH.Del("authorization")

				respB := resp.Body
				// outputting
				fmt.Println("response headers:")
				// fmt.Println(clonedH)
				for header, val := range clonedH {
					fmt.Println("header:", header)
					for _, v := range val {
						fmt.Println("value:", v)
						for range 5 {
							fmt.Print("-")
						}
					}
				}
				for range 10 {
					fmt.Print("-") // seperator for headers and resp body so its easier to read
				}
				body, err := io.ReadAll(respB) // read respB bytes
				respB.Close()                  // close after reading response
				if err != nil {
					continue // skip current iteration if no body
				} else {
					fmt.Println("\nresponse body:")
					for range 10 {
						fmt.Print("-")
					}
					fmt.Println("\n", string(body))
				}

			}
		case "SPAM":
			conf := InitSpamConfig()

			ok := ExtractSpamUrl(parts, store) // conf.reqUrl would hold url
			if !ok {
				continue
			}

			// needed defaults
			req, reqErr := http.NewRequest(http.MethodGet, conf.reqUrl, nil)
			conf.reqCap = 5 // N amount of requests that are to be sent

			if len(parts) <= 2 {
				fmt.Println("please include a request method (-x [method])")
				continue
			}
			uc := strings.ToUpper(parts[2])
			if uc == "-X" {

				conf.reqMethod = parts[3]

				if strings.ToUpper(conf.reqMethod) != "GET" && strings.ToUpper(conf.reqMethod) != "POST" {
					fmt.Println("feature is only compatible with GET/POST requests currently.")
					continue
				}

				// parsing remaining parts after parts[3]
				remainingArgs := parts[4:]

				// get req
				if strings.ToUpper(conf.reqMethod) == "GET" {
					if ok := HandleSpamGet(remainingArgs, conf); !ok {
						continue
					}
				}

				// post req
				if strings.ToUpper(conf.reqMethod) == "POST" {
					ok := HandleSpamPost(remainingArgs, conf)
					if !ok {
						continue
					}
				}

				// initialising final req
				req, reqErr = http.NewRequest(conf.reqMethod, conf.reqUrl, conf.reqBody)
				if reqErr != nil {
					fmt.Println("error initialising request.")
					continue
				}

				res := make(chan *Result) // one channel for all workers
				for i := range conf.reqCap {
					go func(id int) {
						NewWorker(req, id, &conf.client, res)
					}(i)
				}

				for range conf.reqCap {
					r := <-res // receiving results

					if r.Error != nil {
						r.Err()
						continue
					}

					// body
					r.DisplayBody()
				}
			}

		case "FETCH":
			ok := HandleFetch(parts, store)
			if !ok { //
				continue
			}

		case "HELP":
			DisplayCmds()
			continue
		// for exiting
		case "EXIT":
			fmt.Println("exiting will remove saved variables.")
			return

		default:
			fmt.Println("not a valid command.")
			continue
		}
	}
}
