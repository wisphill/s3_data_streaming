package notifier

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	url := "http://localhost:8082"
	fmt.Println("URL:>", url)
	// keepAliveTimeout:= 600 * time.Second
	timeout := 1 * time.Second

	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		TLSHandshakeTimeout: 0 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	makeOneBillionReqs(url, client)
}

func makeOneBillionReqs(url string, client *http.Client) {
	for i := 0; i < 100000; i++ {
		for j := 0; j < 10000; j++ {
			makeRequest(url, client, int64(i*10000+j))
		}
	}
}

func makeRequest(url string, client *http.Client, counter int64) {
	var randomNumber = rand.Intn(10)
	var x = strconv.FormatInt(int64(counter), 10)
	var y = strconv.FormatInt(int64(randomNumber), 10)
	var now = strconv.FormatInt(int64(time.Now().Unix()), 10)

	var jsonStr = []byte(` {"text": "hello world", "content_id": ` + x + `, "client_id": ` + y + `, "timestamp": ` + now + `}`)
	fmt.Println(counter)
	var jsonBytes = bytes.NewReader(jsonStr)
	stringReadCloser := ioutil.NopCloser(jsonBytes)
	defer stringReadCloser.Close()

	req, err := http.NewRequest("POST", url, jsonBytes)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// handler send request
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Errorf("Error")
			return
		}

		defer resp.Body.Close()

		if resp != nil {
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
			return
		}

		_, err = io.Copy(ioutil.Discard, resp.Body) // make sure to read body
		if err != nil {
			return
		}
	}()
}
