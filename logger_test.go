package main

import (
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
		Reqclient:    &client,
		Req:          cl,
		ReqBodySize:  IntPtr(0),
		AllowReqBody: BoolPtr(),
	}

	logConf2 := &LoggerConf{}

	logConf2 = Log(logConf)

	if logConf2.Reqerr != nil {
		t.Fatal("found err:", logConf2.Reqerr.Error())
	}
}
