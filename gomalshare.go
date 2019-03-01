package gomalshare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// LimitKey struct for unmarshal limits of API
type LimitKey struct {
	Limit     string `json:"limit,omitempty"`
	Remaining string `json:"remaining,omitempty"`
}

// HashList struct for unmarshal general hash fields
type HashList struct {
	Md5    string `json:"md5,omitempty"`
	Sha1   string `json:"sha1,omitempty"`
	Sha256 string `json:"sha256,omitempty"`
}

// FileDetails unmarshal special fields
type FileDetails struct {
	HashList
	Ssdeep  string   `json:"ssdeep,omitempty"`
	FType   string   `json:"f_type,omitempty"`
	Sources []string `json:"sources,omitempty"`
}

// Client main struct
type Client struct {
	apiKey string
	url    string
	conn   *http.Client
}

// SearchDetails return searching result
type SearchDetails struct {
	HashList
	TypeSample string `json:"type,omitempty"`
	Added      uint64 `json:"added,omitempty"`
	Source     string `json:"source,omitempty"`
	YaraHits   struct {
		Yara []string `json:"yara,omitempty"`
	} `json:"yarahits,omitempty"`
	Parentfiles []interface{} `json:"parentfiles,omitempty"`
	Subfiles    []interface{} `json:"subfiles,omitempty"`
}

const defaultURL = "http://www.malshare.com/"

// New constructor function
func New(apiKey string, url string) (*Client, error) {
	client := &http.Client{}
	if url == "" {
		url = defaultURL
	}
	if apiKey == "" {
		return nil, errors.New("didn't find API key")
	}
	return &Client{
		conn:   client,
		apiKey: apiKey,
		url:    url,
	}, nil
}
func (c *Client) query(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetSearchResult return details form search sample hashes, sources and file names
func (c *Client) GetSearchResult(str string) (*[]SearchDetails, error) {
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=search&query=%s", c.url, c.apiKey, str)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewBuffer(body))
	var searches []SearchDetails
	for {
		var s SearchDetails
		if err := dec.Decode(&s); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		searches = append(searches, s)

	}
	return &searches, nil
}

// GetLimitKey return allocated number of API key requests per day and remaining
func (c *Client) GetLimitKey() (*LimitKey, error) {
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=getlimit", c.url, c.apiKey)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	var limit LimitKey
	err = json.Unmarshal(body, &limit)
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

// GetListOfTypesFile24 return list of file types & count from the past 24 hours
func (c *Client) GetListOfTypesFile24() (map[string]uint64, error) {
	var filetypes map[string]uint64
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=gettypes", c.url, c.apiKey)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &filetypes)
	if err != nil {
		return nil, err
	}
	return filetypes, nil
}

// GetStoredFileDetails return stored file details
func (c *Client) GetStoredFileDetails(hash string) (*FileDetails, error) {
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=details&hash=%s", c.url, c.apiKey, hash)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	var details FileDetails
	err = json.Unmarshal(body, &details)
	if err != nil {
		return nil, err
	}
	return &details, nil
}

// GetListOfHash24Type return list MD5/SHA1/SHA256 hashes of a specific type from the past 24 hours
func (c *Client) GetListOfHash24Type(typeFile string) (*[]HashList, error) {
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=type&type=%s", c.url, c.apiKey, typeFile)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	var hashes []HashList
	err = json.Unmarshal(body, &hashes)
	if err != nil {
		return nil, err
	}
	return &hashes, nil
}

// GetListOfSource24 return list of sample sources from the past 24 hours
func (c *Client) GetListOfSource24() (*[]string, error) {
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=getsources", c.url, c.apiKey)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	var sources []string
	err = json.Unmarshal(body, &sources)
	if err != nil {
		return nil, err
	}
	return &sources, nil
}

// GetListOfHash24 return list hashes from the past 24 hours
func (c *Client) GetListOfHash24() (*[]HashList, error) {
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=getlist", c.url, c.apiKey)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	var hashes []HashList
	err = json.Unmarshal(body, &hashes)
	if err != nil {
		return nil, err
	}
	return &hashes, nil
}

// DownloadFileFromHash return file for specific hash
func (c *Client) DownloadFileFromHash(hash string) ([]byte, error) {
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=getfile&hash=%s", c.url, c.apiKey, hash)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// UploadFile used for upload using FormData field "upload"
func (c *Client) UploadFile(filename string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("upload", filename)
	if err != nil {
		log.Println("error writing to buffer")
		return err
	}
	fh, err := os.Open(filename)
	if err != nil {
		log.Println("error opening file")
		return err
	}
	defer fh.Close()
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	bodyWriter.Close()
	url := fmt.Sprintf("%sapi.php?api_key=%s&action=upload", c.url, c.apiKey)
	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	resp, err := c.conn.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
