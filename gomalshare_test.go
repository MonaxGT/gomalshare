package gomalshare

import (
	"testing"
	"flag"
	"fmt"
	"os"
)

var apiKey string
var URL string

func init() {
	flag.StringVar(&apiKey,"api", "", "API key MalShare")
	flag.StringVar(&URL,"url", "", "URL MalShare")
	flag.Parse()
	if apiKey == "" {
		fmt.Println("API key is required to run the tests agains MalShare")
		os.Exit(1)
	}

}

func TestGetStoredFileDetails(t *testing.T) {
	conf, err := New(apiKey, URL) // Initiate new connection to API
	if err != nil {
		t.Fatal(err)
	}
	var testMd5 = "95bc3d64f49b03749427fcd6601fa8a7"
	report, err := conf.GetStoredFileDetails(testMd5)
	if err != nil {
		t.Error("error requesting report: ", err)
		return
	}

	if report.Md5 != testMd5 {
		t.Error("requested MD5 doesn't match result: ", testMd5, " - ", report.Md5)
		return
	}
}

func TestGetSearchResult(t *testing.T) {
	conf, err := New(apiKey, URL) // Initiate new connection to API
	if err != nil {
		t.Fatal(err)
	}
	var testMd5 = "95bc3d64f49b03749427fcd6601fa8a7"
	report, err := conf.GetSearchResult(testMd5)
	if err != nil {
		t.Error("error requesting report: ", err)
		return
	}
	for _,v := range *report {
		if v.Md5 == testMd5 {
			return
		}
	}
	t.Error("requested MD5 doesn't match result: ",testMd5)
}

func TestGetLimitKey(t *testing.T) {
	conf, err := New(apiKey, URL) // Initiate new connection to API
	if err != nil {
		t.Fatal(err)
	}
	report,err := conf.GetLimitKey()
	if err != nil {
		t.Error("error requesting report: ", err)
		return
	}
	if report.Limit != "" {
		return
	}
	t.Error("wrong response from API service")
}

func TestDownloadFileFromHash (t *testing.T) {
	conf, err := New(apiKey, URL) // Initiate new connection to API
	if err != nil {
		t.Fatal(err)
	}
	var testMd5 = "95bc3d64f49b03749427fcd6601fa8a7"
	body,err := conf.DownloadFileFromHash(testMd5)
	if err != nil {
		t.Error("error requesting report: ", err)
		return
	}
	if body == nil {
		t.Error("service return nil data")
		return
	}
}

func TestGetListOfHash24(t *testing.T) {
	conf, err := New(apiKey, URL) // Initiate new connection to API
	if err != nil {
		t.Fatal(err)
	}
	report,err := conf.GetListOfHash24()
	if err != nil {
		t.Error("error requesting report: ", err)
		return
	}
	if report == nil {
		t.Error("service return nil data")
		return
	}
}

func TestGetListOfHash24Type(t *testing.T) {
	conf, err := New(apiKey, URL) // Initiate new connection to API
	if err != nil {
		t.Fatal(err)
	}
	testType := "PE32"
	report,err := conf.GetListOfHash24Type(testType)
	if err != nil {
		t.Error("error requesting report: ", err)
		return
	}
	if report == nil {
		t.Error("service return nil data")
		return
	}
}

func TestGetListOfSource24(t *testing.T) {
	conf, err := New(apiKey, URL) // Initiate new connection to API
	if err != nil {
		t.Fatal(err)
	}
	report,err := conf.GetListOfSource24()
	if err != nil {
		t.Error("error requesting report: ", err)
		return
	}
	if report == nil {
		t.Error("service return nil data")
		return
	}
}

func TestGetListOfTypesFile24(t *testing.T) {
	conf, err := New(apiKey, URL) // Initiate new connection to API
	if err != nil {
		t.Fatal(err)
	}
	report,err := conf.GetListOfTypesFile24()
	if err != nil {
		t.Error("error requesting report: ", err)
		return
	}
	if report == nil {
		t.Error("service return nil data")
		return
	}
}