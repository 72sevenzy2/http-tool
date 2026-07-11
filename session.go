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
			if len(parts) < 3 || len(parts) > 3 { // validate arguments
				fmt.Println("please consider the correct format: var <KeyName> <KeyValue>")
			}

			err := store.Set(parts[1], parts[2])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("successful.")
			continue
		case "GET":
			if len(parts) < 2 || len(parts) > 2 {
				fmt.Println("please use the correct format: get <KeyName>")
				continue
			}

			val, err := store.Get(parts[1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println(val)

		case "DEL":
			if len(parts) < 2 || len(parts) > 2 { // validate if parts[1] (key) only is included in input
				fmt.Println("please use the correct format: del <KeyName>")
				continue
			}

			err := store.Del(parts[1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Println("deleted key.")
			continue

			// actual api testing logic (GET only for now)
		case "TEST":
			if len(parts) < 2 { // validate arguments before continuing (other it will panic)
				fmt.Println("variable as second argument does not exit, consider setting a var.")
				continue
			} else {
				val, err := store.Get(parts[1]) // check if parts[1] exists as a var first
				if err != nil {
					fmt.Println(err.Error())
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
						// check if header values exist
						if i+1 >= len(parts) {
							fmt.Println("missing header values.")
							pass = false
							continue
						}

						reqHeaders = strings.SplitN(parts[i+1], ":", 2) // splits headers as:
						// for example: header1:value, splitN() would split it so:
						// map[string]string{
						// 	"header1",
						//  "value",
						// }

						if len(reqHeaders) < 2 || len(reqHeaders) > 2 { // validate length of headers or it will panic during execution
							fmt.Println("please include both header name and value.")
							pass = false
							break
						}
						continue // skip

						// refactored above ^ (keeping this block for future reference)
						// if reqHeaders[0] == "" && reqHeaders[1] == "" {
						// 	// cl.Header.Add(headers[0], headers[1])
						// 	continue
						// }
					}

					// to check if argument is for req methods
					if upc == "-X" {

						// validate if values exist (GET or POST)
						if i+1 >= len(parts) {
							fmt.Println("missing argument values, consider either POST or GET.")
							pass = false
							continue
						}

						newReq := strings.ToUpper(parts[i+1])

						if newReq == "POST" {
							reqType = http.MethodPost

							// validate if request body flag also exists aswell as its data
							if i+2 >= len(parts) {
								fmt.Println("missing request body flag, consider: -d [data]")
								pass = false
								continue
							}

							uppercased1 := strings.ToUpper(parts[i+2]) // normalize parts[i+2] to uppercase upon input

							if uppercased1 == "-D" { // for normal json data

								// validate if values exist
								if i+3 >= len(parts) {
									fmt.Println("please include actual data in json format.")
									pass = false
									continue
								}

								// collect all input after parts[i+3]
								// whilst also removing "\"
								cleaned := strings.ReplaceAll(strings.ReplaceAll(strings.Join(parts[i+3:], ""), " ", ""), "\\", "")

								// assign jsonData to cleaned
								jsonData = cleaned

							}

							if uppercased1 == "-F" { // form mode
								if i+3 >= len(parts) {
									fmt.Println("please include the necessary form data with format: 'title:value'")
									pass = false
									continue
								}

								formParts := strings.SplitN(parts[i+3], ":", 2) // needs to be in a format like so:
								// "title:value", and strings.SplitN(...) would return:
								// map[string]string{
								//		"title": "value",
								// }

								// validate if formParts is of correct length now (or will panic)
								if len(formParts) < 2 || len(formParts) > 2 {
									fmt.Println("please use the correct format.")
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
							if newReq == "GET" {
								reqType = http.MethodGet
							} else {
								fmt.Println("invalid method type.")
								continue
							}
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
					// attaching headers (only if json body data is not included)
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

		case "FETCH":
			ok := HandleFetch(parts, store)
			if !ok { //
				continue
			}

		case "HELP":
			fmt.Println("usage:")
			fmt.Println("var <VarName> <Value>")
			fmt.Println("to retrieve values:")
			fmt.Println("get <VarName>")
			fmt.Println("to test API's:")
			fmt.Println("test <VarName> -H content-type:application -X post -D { 'name': 'name' } or -X post -F <formTitle:formValue>")
			fmt.Println("though flags shown above are optional, but make sure <VarName> is a valid url.")
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
