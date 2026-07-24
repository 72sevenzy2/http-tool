package main

import (
	"net/http"
	"testing"
)

func boolre() *bool {
	b := true
	return &b
}

func intPr(i int) *int { return &i }

func TestLogger(t *testing.T) {
	cl, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	logConf := &LoggerParamConf{
		Reqclient:    &client,
		Req:          cl,
		ReqBodySize:  intPr(0),
		AllowReqBody: boolre(),
	}

	logConf2 := &LoggerConf{}

	logConf2 = Log(logConf)

	if logConf2.Reqerr != nil {
		t.Fatal("found err:", logConf2.Reqerr.Error())
	}
}
