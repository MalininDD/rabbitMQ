package internal

import (
	"fmt"
	"io/ioutil"
	"mvdan.cc/xurls/v2"
	"net/http"
)

func Grabber(url string) []string {
	req, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	rxStrict := xurls.Strict()
	return rxStrict.FindAllString(string(b), -1)
}
