package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// fetch cmd
func HandleFetch(parts []string, store *Data) bool {
	if len(parts) > 1 {
		fmt.Println("invalid fetch format: only run <FETCH> to retrieve all set variables.")
		return false
	}

	vals, ok := store.GetAll()
	if !ok {
		fmt.Println("no variables set.")
		return false
	}
	for v, i := range vals {
		fmt.Println(i, v)
	}

	return true
}

// var cmd
func HandleVar(parts []string, store *Data) bool {
	if len(parts) < 3 || len(parts) > 3 { // validate arguments
		fmt.Println("please consider the correct format: var <KeyName> <KeyValue>")
		return false
	}

	err := store.Set(parts[1], parts[2])
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("successful.")
	return true
}

// get cmd
func HandleGet(parts []string, store *Data) bool {
	if len(parts) < 2 || len(parts) > 2 {
		fmt.Println("please use the correct format: get <KeyName>")
		return false
	}

	val, err := store.Get(parts[1])
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println(val)
	return true
}

// delete cmd
func HandleDel(parts []string, store *Data) bool {
	if len(parts) < 2 || len(parts) > 2 { // validate if parts[1] (key) only is included in input
		fmt.Println("please use the correct format: del <KeyName>")
		return false
	}

	msg, ok := store.Del(parts[1])
	if !ok {
		fmt.Println(msg)
		return false
	}

	fmt.Println("deleted key.")
	return true
}

// modules for test cmd:

func HandleTestValdiation(parts []string, store *Data) (string, bool) {
	if len(parts) < 2 {
		fmt.Println("variable as second argument does not exit, consider setting a var.")
		return "", false
	}

	// check if variable stored exists
	val, err := store.Get(parts[1])
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}

	return val, true
}

func HandleHeaderAppending(parts []string, partsIndex int, uppercasedH string) ([]string, bool) {
	// validate if values exist
	if partsIndex+1 >= len(parts) {
		fmt.Println("missing header values.")
		return nil, false
	}

	var returnSlice []string
	returnSlice = strings.SplitN(parts[partsIndex+1], ":", 2) // would create 2 substrings on : seperation

	if len(returnSlice) < 2 || len(returnSlice) > 2 {
		fmt.Println("please include both header name and value.")
		return nil, false
	}

	return returnSlice, true
}

func HandleRequestMethodValidation(parts []string, partsIndex int) (string, bool) {
	if partsIndex+1 >= len(parts) { // check if arguments after partsIndex + 1 exist
		fmt.Println("missing argument values, consider either POST or GET.")
		return "", false
	}

	return strings.ToUpper(parts[partsIndex+1]), true
}

func HandlePostValidation(parts []string, partsIndex int) (string, bool) {
	if partsIndex+2 >= len(parts) {
		fmt.Println("missing request body flag, consider: -d [data]")
		return "", false
	}

	return strings.ToUpper(parts[partsIndex+2]), true
}

// small util func for validating json inputs
func IsValidJson(str []byte) bool {
	if v := json.Valid(str); !v {
		return false
	}
	return true
}

func HandleJsonDataInput(parts []string, partsIndex int) (string, bool) {
	if partsIndex+3 >= len(parts) {
		fmt.Println("please include actual data in json format.")
		return "", false
	}

	cleaned := strings.Join(parts[partsIndex+3:], " ") // join json inputs

	valid := IsValidJson([]byte(cleaned))
	if !valid {
		fmt.Println("invalid json format.")
		return "", false
	}

	return cleaned, true
}

func HandleFormDataInput(parts []string, partsIndex int) ([]string, bool) {
	if partsIndex+3 >= len(parts) {
		fmt.Println("please include the necessary form data with format: 'title:value'")
		return nil, false
	}

	formParts := strings.SplitN(parts[partsIndex+3], ":", 2) // splits form data seperated by ":", into substrings

	// validate if formParts is of correct length before accessing slice (or will panic)
	if len(formParts) < 2 || len(formParts) > 2 {
		fmt.Println("please use the correct format.")
		return nil, false
	}

	return formParts, true
}

// spam cmd worker
func NewWorker(req *http.Request, id int, client *http.Client) (string, bool) {
	_, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("worker number %d had failed.", id), false
	}

	return fmt.Sprint("worker finished with id:", id), true
}

// help cmd func

func DisplayCmds() {
	fmt.Println("usage:")
	fmt.Println("var <VarName> <Value>")
	fmt.Println("to retrieve values:")
	fmt.Println("get <VarName>")
	fmt.Println("to test API's:")
	fmt.Println("test <VarName> -H content-type:application -X post -D { 'name': 'name' } or -X post -F <formTitle:formValue>")
	fmt.Println("though flags shown above are optional, but make sure <VarName> is a valid url.")
}
