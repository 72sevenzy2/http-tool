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