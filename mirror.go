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

func CheckReflectedParameters(url_t string, parameters url.Values, wg *sync.WaitGroup, sem chan bool) {
	defer wg.Done()
	<-sem
	payloads := url.Values{}
	i := 0
	reversed_payload := make(map[string]string)
	for name := range parameters {
		payload := "swagonlolnow" + fmt.Sprint(i)
		payloads.Set(name, payload)
		reversed_payload[payload] = name
		i++
	}
	var new_url string
	encoded_payloads := payloads.Encode()
	if CheckContains(url_t) {
		new_url = fmt.Sprintf("%s&%s", url_t, encoded_payloads)
	} else {
		new_url = fmt.Sprintf("%s?%s", url_t, encoded_payloads)
	}
	resp, err := http.Get(new_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	bodyString := string(bodyBytes)
	var matches []string
	for x := 0; x <= i; x++ {
		pattern := "swagonlolnow" + fmt.Sprint(x)
		re, _ := regexp.Compile(pattern)
		matches = append(matches, re.FindString(bodyString))
	}
	reflected_url := url.Values{}
	for payload, name := range reversed_payload {
		for _, match := range matches {
			if match == payload {
				reflected_url.Set(name, "1")
			}
		}
	}
	if len(reflected_url) > 0 {
		encoded_reflected_payloads := reflected_url.Encode()
		new_url = strings.Split(url_t, "?")[0]
		new_url = fmt.Sprintf("%s?%s", new_url, encoded_reflected_payloads)
		fmt.Println(new_url)
	} else {
		return
	}
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
