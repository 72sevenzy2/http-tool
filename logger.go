package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// return config
type LoggerConf struct {
	Req     *http.Response
	Reqtime time.Duration
	Reqbody string
	Reqlen  *int
	Reqerr  error
}

// small utils for LoggerConf

func (b LoggerConf) DisplayReqStatus() string {
	return b.Req.Status
}

func (b LoggerConf) DisplayReqLatency() {
	fmt.Println("request latency:", b.Reqtime)
}

func (b LoggerConf) DisplayReqBodyLen() {
	fmt.Println("request body length:", b.Reqlen)
}

func (b LoggerConf) DisplayReqBody() {
	fmt.Println(b.Reqbody)
}

func (b LoggerConf) Err() {
	fmt.Println(b.Reqerr.Error())
}


// for params
type LoggerParamConf struct {
	Reqclient    *http.Client
	Req          *http.Request
	AllowReqBody *bool
	ReqBodySize  *int
}

func Log(conf *LoggerParamConf) *LoggerConf {
	start := time.Now()
	resp, err := conf.Reqclient.Do(conf.Req)
	if err != nil {
		return &LoggerConf{
			Req:     nil,
			Reqtime: 0,
			Reqbody: "",
			Reqlen:  IntPtr(0),
			Reqerr:  err,
	 	}
	}
	end := time.Since(start)

	// fmt.Println("visited to:", req.URL.Path)
	// fmt.Println("method:", req.Method)
	fmt.Printf("visited to %s, with method %s", conf.Req.URL.Path, conf.Req.Method)
	fmt.Println(" full url:", conf.Req.URL)
	// request query
	if conf.Req.URL.RawQuery != "" {
		fmt.Println("with query:", conf.Req.URL.RawQuery)
	}

	// exlude sensitive headers
	fmt.Println("request header details:")
	newHeaders := conf.Req.Header.Clone()
	newHeaders.Del("Authorization")
	fmt.Println(newHeaders) // then display

	// request body printing
	var bodyprev string  // body preview (string) var
	var max *int         // body size var
	var bodybytes []byte // body size (in bytes)

	if *conf.AllowReqBody {
		max = conf.ReqBodySize
		bodybytes, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}

		resp.Body = io.NopCloser(bytes.NewBuffer(bodybytes))

		fmt.Println("max:", *max, "len:", len(bodybytes))
		if *max > 0 && len(bodybytes) > *max {
			bodybytes = bodybytes[:*max]
		}
		fmt.Println("after truncating:", len(bodybytes))
		bodyprev = string(bodybytes)

	}

	final := len(bodybytes)

	return &LoggerConf{
		Req:     resp,
		Reqtime: end,
		Reqbody: bodyprev,
		Reqlen:  &final,
		Reqerr:  nil,
	}
}
