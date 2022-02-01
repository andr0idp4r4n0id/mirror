package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func PrintReflectedUrls(url_t string) {
	fmt.Println(url_t)
}

func CheckReflection(resp *http.Response) bool {
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	bodyString := string(body_bytes)
	return strings.Contains(bodyString, "swagonswagggyswag")
}

func SendHttpRequest(url_t string) *http.Response {
	resp, err := http.Get(url_t)
	if err != nil {
		return nil
	}
	return resp
}

func ReplaceValuesInUlr(url_t string, query url.Values) string {
	url_t = strings.Split(url_t, "?")[0]
	url_t += "?" + query.Encode()
	return url_t

}

func ReplaceValuesInQuery(query url.Values) url.Values {
	for param := range query {
		query[param][0] = "swagonswagggyswag"
	}
	return query
}

func CheckReflectedParameters(url_t string, query url.Values) {
	query = ReplaceValuesInQuery(query)
	url_t = ReplaceValuesInUlr(url_t, query)
	resp := SendHttpRequest(url_t)
	if resp == nil {
		return
	}
	if !CheckReflection(resp) {
		return
	}
	PrintReflectedUrls(url_t)
	defer resp.Body.Close()
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	conc := flag.Int("concurrency", 10, "concurrency level.")
	flag.Parse()
	var wg sync.WaitGroup
	for i := 0; i < *conc; i++ {
		for reader.Scan() {
			url_t := reader.Text()
			uri, _ := url.Parse(url_t)
			wg.Add(1)
			go func() {
				CheckReflectedParameters(url_t, uri.Query())
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
