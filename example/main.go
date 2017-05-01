package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fd0/drainchecker"
)

var readBody = flag.Bool("readbody", false, "read HTTP response body")

func init() {
	flag.Parse()
}

func main() {
	http.DefaultTransport = drainchecker.RoundTripper(http.DefaultTransport)
	http.DefaultClient.Transport = http.DefaultTransport

	for i := 0; i < 10; i++ {
		res, err := http.Get("https://www.google.com")
		if err != nil {
			panic(err)
		}

		fmt.Printf("status: %v\n", res.Status)

		if *readBody {
			_, err = io.Copy(ioutil.Discard, res.Body)
			if err != nil {
				panic(err)
			}
		}

		err = res.Body.Close()
		if err != nil {
			panic(err)
		}
	}
}
