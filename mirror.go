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
	return strings.Contains(bodyString, "swagonswagggyswag123")
}

func SendHttpRequest(url_t string) *http.Response {
	resp, err := http.Get(url_t)
	if err != nil {
		return nil
	}
	return resp
}

func ReplaceValuesInUlr(url_t string, query string) string {
	url_t += "?" + query
	return url_t

}

func ReplaceValuesInQuery(param, query string) string {
	query_strings := url.Values{}
	query_strings.Set(param, query)
	query_strings_encoded := query_strings.Encode()
	return query_strings_encoded
}

func CheckReflectedParameters(url_t string, param string, query string) {
	query_strings_encoded := ReplaceValuesInQuery(param, query)
	url_t = ReplaceValuesInUlr(url_t, query_strings_encoded)
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
			url_t = strings.Split(url_t, "?")[0]
			for key := range uri.Query() {
				param := key
				wg.Add(1)
				go func() {
					CheckReflectedParameters(url_t, param, "swagonswagggyswag123")
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}
}
