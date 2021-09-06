package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

func CheckContains(url_t string) bool {
	re := regexp.MustCompile(`\?\w.+`)
	matched := re.MatchString(url_t)
	return matched
}

func SetPayloads(parameters url.Values, reversed_payload map[string]string, payloads url.Values, i int) {
	for name := range parameters {
		payload := "swagonlolnow" + fmt.Sprint(i)
		payloads.Set(name, payload)
		reversed_payload[payload] = name
		i++
	}
}

func EncodePayloads(payloads url.Values) string {
	return payloads.Encode()
}

func SendHttpRequestReadResponseBody(new_url string) []byte {
	resp, err := http.Get(new_url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return bodyBytes
}

func CheckMatches(bodyString string, i int) []string {
	var matches []string
	for x := 0; x <= i; x++ {
		pattern := "swagonlolnow" + fmt.Sprint(x)
		re, _ := regexp.Compile(pattern)
		matches = append(matches, re.FindString(bodyString))
	}
	return matches
}

func FindPayloadInReversePayloads(matches []string, reversed_payload map[string]string) url.Values {
	reflected_url := url.Values{}
	for payload, name := range reversed_payload {
		for _, match := range matches {
			if match == payload {
				reflected_url.Set(name, "1")
			}
		}
	}
	return reflected_url

}

func PrintReflections(reflected_url url.Values, new_url string, url_t string) {
	if len(reflected_url) > 0 {
		encoded_reflected_payloads := reflected_url.Encode()
		url := strings.Split(url_t, "?")[0]
		url = fmt.Sprintf("%s?%s", url, encoded_reflected_payloads)
		fmt.Println(url)
	} else {
		return
	}
}

func CheckReflectedParameters(url_t string, parameters url.Values, wg *sync.WaitGroup, sem chan bool) {
	defer wg.Done()
	<-sem
	i := 0
	payloads := url.Values{}
	reversed_payload := make(map[string]string)
	SetPayloads(parameters, reversed_payload, payloads, i)
	encoded_payloads := EncodePayloads(payloads)
	var new_url string
	if CheckContains(url_t) {
		new_url = fmt.Sprintf("%s&%s", url_t, encoded_payloads)
	} else {
		new_url = fmt.Sprintf("%s?%s", url_t, encoded_payloads)
	}
	bodyBytes := SendHttpRequestReadResponseBody(new_url)
	if bodyBytes != nil {
		return
	}
	bodyString := string(bodyBytes)
	matches := CheckMatches(bodyString, i)
	reflected_url := FindPayloadInReversePayloads(matches, reversed_payload)
	PrintReflections(reflected_url, new_url, url_t)
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	conc := flag.Int("concurrency", 10, "concurrency level.")
	sem := make(chan bool, *conc)
	var wg sync.WaitGroup
	for reader.Scan() {
		url_t := reader.Text()
		uri, _ := url.Parse(url_t)
		wg.Add(1)
		sem <- true
		go CheckReflectedParameters(url_t, uri.Query(), &wg, sem)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	wg.Wait()
}
