package main

import (
	"fmt"
	"log"
	"os"

	"bytes"
	"encoding/json"
	"github.com/iris-contrib/color"
	"github.com/itang/godict"
	"github.com/itang/gotang"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	word := doParseWordFromArgs()
	from := godict.Word{W: word, L: godict.LANG_EN}

	t := godict.Translator163{}

	ret, err := t.Translate(from, godict.LANG_CN)
	if err != nil {
		log.Fatalln("出错了", color.RedString(err.Error()))
	}

	fmt.Printf("%s:\n", color.RedString(word))
	fmt.Println(">", color.BlueString(ret))

	fmt.Println("\nPost to cloud...")
	gotang.Time(func() {
		done := make(chan Result)
		timer := time.NewTimer(time.Millisecond * 2000)

		go httpPostAsString("http://dict.godocking.com/api/dict/logs", postRequest{From: word, To: ret}, done)

		select {
		case ret := <-done:
			value := ret.ok_or(func(err error) interface{} {
				return err.Error()
			})
			fmt.Printf(" -> response: %v\n", value)
		case <-timer.C:
			fmt.Println("timeout...")
		}
	})
}

type postRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func doParseWordFromArgs() string {
	if len(os.Args) < 2 {
		log.Fatalln(">", color.RedString("请输入要翻译的词汇"))
	}

	word := os.Args[1]
	return word
}

func httpPostAsString(url string, req postRequest, done chan Result) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(req)
	resp, err := http.Post(url, "application/json; charset=utf-8", &buf)

	if err != nil {
		done <- Result{nil, err}
	} else {
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			done <- Result{nil, err}
		}
		done <- Result{string(content), nil}
	}
}

type Result struct {
	Value interface{}
	Err   error
}

func (ret Result) flat() (interface{}, error) {
	return ret.Value, ret.Err
}

func (ret Result) ok_or(f func(err error) interface{}) interface{} {
	if ret.Err != nil {
		return f(ret.Err)
	}
	return ret.Value
}
