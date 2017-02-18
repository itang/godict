package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/iris-contrib/color"
	"github.com/itang/godict"
	"github.com/itang/gotang"
	"github.com/pkg/errors"
)

func main() {
	word, err := parseWordFromArgs()
	if err != nil {
		fmt.Printf("INFO: %s\n", color.RedString(err.Error()))
		return
	}

	fmt.Printf("%s:\n", color.GreenString(word))
	ret, err := godict.Translator163{}.Translate(godict.Word{W: word, L: godict.LANG_EN}, godict.LANG_CN)
	if err != nil {
		fmt.Printf("ERROR: %s\n", color.RedString(err.Error()))
		return
	}

	fmt.Printf("\t->: %s\n", color.BlueString(ret))

	tryPostToCloud(word, ret)
}

const MAX_TO_CHARS = 100

func tryPostToCloud(from, to string) {
	fmt.Println("\ntry post to cloud...")
	if len(to) > MAX_TO_CHARS {
		fmt.Printf("INFO: Too large content(%v bytes), ignore post.\n", len(to))
		return
	}

	gotang.Time(func() {
		done := make(chan Result)
		timer := time.NewTimer(time.Millisecond * 2000)

		go httpPostAsString("http://dict.godocking.com/api/dict/logs", postRequest{From: from, To: to}, done)

		select {
		case ret := <-done:
			value := ret.okOrElse(func(err error) interface{} {
				return err.Error()
			})
			fmt.Printf("\t->: %v\n", value)
		case <-timer.C:
			fmt.Println("timeout...")
		}
	})
}

type postRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func parseWordFromArgs() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("请输入要翻译的词汇")
	}

	word := os.Args[1]
	return word, nil
}

func httpPostAsString(url string, req postRequest, done chan Result) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(req)
	resp, err := http.Post(url, "application/json; charset=utf-8", &buf)
	if err != nil {
		done <- Result{nil, err}
		return
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		done <- Result{nil, err}
		return
	}

	done <- Result{string(content), nil}
}

// Result type
type Result struct {
	Value interface{}
	Err   error
}

func (ret Result) flat() (interface{}, error) {
	return ret.Value, ret.Err
}

func (ret Result) okOrElse(f func(err error) interface{}) interface{} {
	if ret.Err != nil {
		return f(ret.Err)
	}
	return ret.Value
}
