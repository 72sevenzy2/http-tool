package main

import "fmt"

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

	err := store.Del(parts[1])
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	
	fmt.Println("deleted key.")
	return true
}

// for test cmd

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