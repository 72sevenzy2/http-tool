// to hold user inputs during session mode - storing variables in maps to then be retrieved later using dynamic lookups.

package main

import (
	"errors"
	"fmt"
)

type Data struct {
	data_storage map[string]string
}

// new db
func NewStore() *Data {
	return &Data{
		make(map[string]string),
	}
}

// utility get/set functions for data map:

// for both strings and ints
func (d *Data) Get(keyname any) (string, error) {
	newKey := Normalize(keyname)

	val, ok := d.data_storage[newKey]
	if !ok {
		return "", errors.New("key does not exist.")
	}

	return val, nil
}

func (d *Data) Set(keyname any, value string) error {
	if len(value) == 1 {
		return errors.New("please include a value thats over 1 character.")
	}

	newK := Normalize(keyname)

	d.data_storage[newK] = value
	return nil
}

// del func
func (d *Data) Del(keyname any) (string, bool) {
	newk := Normalize(keyname)

	if _, v := d.data_storage[newk]; !v {
		return "key does not exist.", false
	}

	delete(d.data_storage, newk)
	return fmt.Sprintf("deleted key with value: %s", d.data_storage[newk]), true
}

// func to get all values from data_storage
func (d *Data) GetAll() (map[string]string, bool) {
	res := make(map[string]string, len(d.data_storage)) // initialise the size as number of elements in data_storage to reduce size allocated for this map

	if len(d.data_storage) != 0 {
		for v, i := range d.data_storage {
			res[i] = v
			return res, true
		}
	}
	return nil, false
}
