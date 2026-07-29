package main

import (
	"log"
	"net/http"
	"testing"
)

func BoolPtr() *bool {
	b := true
	return &b
}

func TestLogger(t *testing.T) {
	cl, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	logConf := &LoggerParamConf{
		ReqNo:        IntPtr(8),
		Reqclient:    &client,
		Req:          cl,
		ReqBodySize:  IntPtr(0),
		AllowReqBody: BoolPtr(),
	}

	logConf2 := &LoggerConf{}

	logConf2 = Log(logConf)

	for range len(logConf2.Results) {
		log.Fatal("found error:")
		logConf2.Err()
	}
}
