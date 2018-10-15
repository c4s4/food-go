package main

import (
	"os"
	"net/http"
	"io/ioutil"
	"strconv"
)

const (
	Min = 1000
	Max = 11000
)

func main() {
	url := os.Args[1]
	for i := Min; i <= Max; i++ {
		response, err := http.Get(url+strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		println(string(content))
	}
}

// golang: 4,868 s
// newapi: 3 min 55,883 s = 235 s
