package main

import (
	"fmt"
	"os"

	"net/http"
	"net/url"

	"github.com/op/go-logging"
	"golang.org/x/net/html"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	distURL = kingpin.Arg("url", "url to parse").Required().String()
)

var log = logging.MustGetLogger("example")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	fmt.Printf("----- Parse URL : %s\n", *distURL)

	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)

	_, err := url.Parse(*distURL)
	if err != nil {
		log.Fatal(err)
	}

	resp, _ := http.Get(*distURL)
	//bytes, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println("HTML:\n\n", string(bytes))
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				fmt.Printf("We found a link!\n")
				for _, a := range t.Attr {
					if a.Key == "href" {
						log.Info("Found href:", a.Val)
						break
					}
				}
			}
		}
	}
	resp.Body.Close()
}
