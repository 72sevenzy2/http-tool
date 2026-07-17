package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
			var (
				pass    bool // determine whether final output block runs
				reqType string
				// bodyData map[string]string
				jsonData string
			)
			pass = true

			// request dependent variables
			var (
				cl    *http.Request
				clErr error

				// for appending headers
				reqHeaders []string // reqHeaders being globally accessed by this scope, so headers can be assigned after initialising request.

				formBody io.Reader // to hold form data
			)

			// default request (GET) so cl does not stay nil if user were not to pass -X
			cl, clErr = http.NewRequest(http.MethodGet, val, nil)
			// clErr is handled below, after flags are parsed and appropriate requests are made.

			// parse header arguments
			for i := range parts {
				upc := strings.ToUpper(parts[i]) // normalize "-h" to all uppercase

				if upc == "-H" {
					headerSlice, ok := HandleHeaderAppending(parts, i, upc)
					if !ok {
						pass = false
						continue
					}

					reqHeaders = headerSlice // store headerSLice contents which were extracted.
				}

				// to check if argument is for req methods
				if upc == "-X" {
					newReqMethod, exists := HandleRequestMethodValidation(parts, i)
					if !exists {
						pass = false
						continue
					}

					if newReqMethod == "POST" {
						reqType = http.MethodPost

						uppercased1, exists2 := HandlePostValidation(parts, i)
						if !exists2 {
							pass = false
							continue
						}

						if uppercased1 == "-D" { // for normal json data

							cleanedData, ok := HandleJsonDataInput(parts, i)
							if !ok {
								pass = false
								continue
							}

							// assign jsonData to cleaned
							jsonData = cleanedData

						}

						if uppercased1 == "-F" { // form mode
							formParts, ok := HandleFormDataInput(parts, i)
							if !ok {
								pass = false
								continue
							}

							// appending the form values
							data := url.Values{}

							data.Set(formParts[0], formParts[1])
							enc := data.Encode() // encode

							formBody = strings.NewReader(enc)
						}

					} else {
						if newReqMethod == "GET" {
							reqType = http.MethodGet
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
				if jsonData != "" && reqType == http.MethodPost {
					cl, clErr = http.NewRequest(reqType, val, strings.NewReader(jsonData))

					// set content type to json after creating request
					cl.Header.Add("Content-Type", "application/json")
					isJsonH = true
				} else {
					if jsonData == "" && reqType == http.MethodPost { // meaning that its form related request body data
						cl, clErr = http.NewRequest(reqType, val, formBody)

						// setting appropriate header afterwards
						cl.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						isFormH = true
					}
					if reqType == http.MethodGet { // standard get method
						cl, clErr = http.NewRequest(http.MethodGet, val, nil)
					}

				}

				if clErr != nil {
					fmt.Println(clErr)
					continue
				}

				if !pass {
					continue
				}

				if len(reqHeaders) == 2 {
					// attaching headers (only if it json data, or form data isnt included as it would override those headers.)
					if !isJsonH && !isFormH {
						cl.Header.Add(reqHeaders[0], reqHeaders[1])
					}
				}

				resp, err2 := client.Do(cl) // send the request to the url provided
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
			var (
				client http.Client

				reqMethod string
				reqUrl    string
				reqBody   io.Reader
			)

			// needed defaults
			req, reqErr := http.NewRequest(http.MethodGet, reqUrl, nil)
			reqCap := 5 // N amount of requests that are to be sent

			// if i+1 >= len(parts) {
			if len(parts) <= 1 {
				fmt.Println("please include a url.")
				continue
			}

			varVal, err := store.Get(parts[1])
			if err != nil { // it isnt a variable, and is real url.
				reqUrl = parts[1]
			}

			reqUrl = varVal

			if len(parts) <= 2 {
				fmt.Println("please include a request method (-x [method])")
				continue
			}
			uc := strings.ToUpper(parts[2])
			if uc == "-X" {

				reqMethod = parts[3]

				if strings.ToUpper(reqMethod) != "GET" && strings.ToUpper(reqMethod) != "POST" {
					fmt.Println("feature is only compatible with GET/POST requests currently.")
					continue
				}

				// parsing remaining parts after parts[3]
				remainingArgs := parts[4:]

				// get req
				if strings.ToUpper(reqMethod) == "GET" {
					if len(remainingArgs) <= 2 {
						fmt.Println("please include number of requests that is to be sent in double quotes (-C ['N']).")
						continue
					}

					// initialise get req
					req, reqErr = http.NewRequest(http.MethodGet, reqUrl, nil)
					if reqErr != nil {
						fmt.Println("error initialising request.")
						continue
					}

					s := strings.Trim(remainingArgs[2], "\"")
					d, err := strconv.Atoi(s)
					if err != nil {
						fmt.Println("invalid number.")
						continue
					}

					reqCap = d

					// initialise workers
					for i := range reqCap {
						go func() {
							msg, ok := NewWorker(req, i, &client)
							if !ok {
								str := fmt.Sprintf("worker with id %d failed to make request.", i)
								fmt.Println(str)
								return
							}

							fmt.Println(msg)
						}()
					}
				}

				if strings.ToUpper(reqMethod) == "POST" {
					if len(remainingArgs) <= 2 {
						fmt.Println("please include number of requests that is to be sent in double quotes (-C ['N']).")
						continue
					}

					// req initialisation later

					// todo: validate if flag is -C
					// for reqcap and parsing -C value
					s := strings.Trim(remainingArgs[2], "\"")
					d, err := strconv.Atoi(s)
					if err != nil {
						fmt.Println("not a valid number.")
						continue
					}
					reqCap = d

					// body extraction
					if len(remainingArgs) <= 4 {
						fmt.Println("please include a json body (-d [json]).")
						continue
					}

					// todo: validate if flag is -d

					reqBody = strings.NewReader(strings.Join(remainingArgs[4:], " "))
					valid := IsValidJson([]byte(strings.Join(remainingArgs[4:], " ")))
					if !valid {
						fmt.Println("invalid json.")
						continue
					}

					req, reqErr = http.NewRequest(http.MethodPost, reqUrl, reqBody)
					if reqErr != nil {
						fmt.Println("error initialising request.")
						continue
					}

					// start workers
					for id := range reqCap {
						go func() {
							msg, ok := NewWorker(req, id, &client)
							if !ok {
								str := fmt.Sprintf("worker with id %d failed to make request.", id)
								fmt.Println(str)
								return
							}

							fmt.Println(msg)
						}()
					}
				}

				// initialise request (for GET)
				// req, reqErr = http.NewRequest(reqMethod, reqUrl, nil)

				// if strings.ToUpper(reqMethod) == "POST" {
				// 	if len(parts) <= 4 { // for the body
				// 		fmt.Println("missing json data.")
				// 		continue
				// 	}

				// 	cleanedJson := strings.Join(parts[4:], " ")
				// 	// validate json
				// 	if ok := json.Valid([]byte(cleanedJson)); !ok {
				// 		fmt.Println("invalid json.")
				// 		continue
				// 	}

				// 	req, reqErr = http.NewRequest()
				// }

				// // handling request error
				// if reqErr != nil {
				// 	fmt.Println("error initialising request.")
				// 	continue // skip iteration
				// }

				// if len(parts) <= 4 {
				// 	fmt.Println("please include a number of requests to be sent in double quotes (-c ['N amounts of requsts'])")
				// 	continue
				// }
				// if strings.ToUpper(parts[4]) == "-C" {
				// 	if len(parts) <= 5 {
				// 		fmt.Println("please include a number of requests in double quotes.")
				// 		continue
				// 	}
				// 	s := strings.Trim(parts[5], "\"") // remove \
				// 	r, err := strconv.Atoi(s)         // parse to int
				// 	if err != nil {
				// 		fmt.Println("invalid number.")
				// 		continue
				// 	}
				// 	reqCap = r
				// 	client := http.Client{} // workers share the same client details instead of getting their own client struct

				// 	// start workers
				// 	for id := range reqCap {
				// 		go func() {
				// 			msg, ok := NewWorker(req, id, &client)
				// 			if !ok {
				// 				str := fmt.Sprintf("worker with id %d failed to make request.", id)
				// 				fmt.Println(str)
				// 				return
				// 			}
				// 			time.Sleep(time.Second)
				// 			fmt.Println(msg)
				// 		}()
				// 	}
				// }
				// uc2 := strings.ToUpper(parts[i+3])
				// if uc2 == "GET" {
				// 	req, reqErr = http.NewRequest(http.MethodGet, reqUrl, nil)
				// }
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
