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
					ok, headerSlice := HandleHeaderAppending(parts, i, upc)
					if !ok {
						pass = false
						continue
					}

					reqHeaders = headerSlice // store headerSLice contents which were extracted.
				}

				// to check if argument is for req methods
				if upc == "-X" {
					exists, newReqMethod := HandleRequestMethodValidation(parts, i)
					if !exists {
						pass = false
						continue
					}

					if newReqMethod == "POST" {
						reqType = http.MethodPost

						exists2, uppercased1 := HandlePostValidation(parts, i)
						if !exists2 {
							pass = false
							continue
						}

						if uppercased1 == "-D" { // for normal json data

							ok, cleanedData := HandleJsonDataInput(parts, i)
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
