[![GoDoc](https://godoc.org/github.com/MonaxGT/gomalshare?status.svg)](http://godoc.org/github.com/MonaxGT/gomalshare)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/53870d8fa2174624ad09895fb4fc2f1a)](https://app.codacy.com/app/MonaxGT/gomalshare?utm_source=github.com&utm_medium=referral&utm_content=MonaxGT/gomalshare&utm_campaign=Badge_Grade_Dashboard)
[![Build Status](https://travis-ci.com/MonaxGT/gomalshare.svg?branch=master)](https://travis-ci.com/MonaxGT/gomalshare)
[![Go Report Card](https://goreportcard.com/badge/github.com/MonaxGT/gomalshare)](https://goreportcard.com/report/github.com/MonaxGT/gomalshare)

MalShare client library 
---------------------
MalShare is a free Malware repository providing researchers access to samples, malicous feeds, and Yara results. 
Link to Malshare: 

* [github](https://github.com/malshare)
* [official site](http://www.malshare.com)
* [twitter](https://twitter.com/mal_share)

Usage example
------------------------------------------------

```sh
go get -u github.com/MonaxGT/gomalshare
```

```sh
go test -api APIKEY -url URL
```

Simple example using library in cmd/gomalshare/main.go

```go
package main

import (
	"flag"
	"fmt"

	"github.com/MonaxGT/gomalshare"
)

func main() {
	apiKeyPtr := flag.String("api", "", "API key MalShare")
	urlPtr := flag.String("url", "", "URL MalShare")
	flag.Parse()
	var err error
	var conf *gomalshare.Client

	// init function
	conf, err = gomalshare.New(*apiKeyPtr, *urlPtr) // Initiate new connection to API
	if err != nil {
		panic(err)
	}

	// example with return list of hashes last 24 hours
	var list24 *[]gomalshare.HashList
	list24, _ = conf.GetListOfHash24()
	fmt.Println(list24)

	// example with return list of types of downloading files last 24 hours
	typeCount, _ := conf.GetListOfTypesFile24()
	fmt.Println(typeCount)

	// example with return current api key limit
	var limitKey *gomalshare.LimitKey
	limitKey, _ = conf.GetLimitKey()
	fmt.Println(limitKey)

	// example with return information of files by using sample
	var search *[]gomalshare.SearchDetails
	search, err = conf.GetSearchResult("emotet")
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range *search {
		fmt.Println(v.Md5)
	}
	// example upload file
	filename := "test.test"
	err = conf.UploadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	// example for download file by hash request
	file, err := conf.DownloadFileFromHash("95bc3d64f49b03749427fcd6601fa8a7")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(file))
}
```
