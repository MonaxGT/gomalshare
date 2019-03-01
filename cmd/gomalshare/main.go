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
