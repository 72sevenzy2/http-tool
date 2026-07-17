package main

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	store := NewStore()
	er := store.Set(1, "test")
	if er == nil {
		fmt.Println("data set successfully")
	}
	fmt.Println("could not save data")

	_, ok := store.Del(1)
	if !ok {
		fmt.Println("key does no exist.")
	}

	val, err1 := store.Get(1)
	if err1 == nil {
		fmt.Println("data exists:", val)
	}
	fmt.Println("couldnt retrieve data") // expected output for del testing

	err3 := store.Set("test", "hi")
	if err3 != nil {
		fmt.Println("couldnt set data.")
	}

	vals, ok3 := store.GetAll()
	if !ok3 {
		fmt.Println("not vars set.")
	}
	for i, v := range vals {
		fmt.Println(v, i)
	}
}
