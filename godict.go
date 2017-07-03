package godict

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/itang/gotang"
	"github.com/pkg/errors"
)

const (
	LANG_EN = iota
	LANG_CN
)

type Lang int

type Word struct {
	W string
	L Lang
}

type Translator interface {
	Translate(from Word, to Lang) (ret string, err error)
}

type Record interface {
	Record(from Word, to Word)
}

type Translator163 struct {
}

func (t Translator163) Translate(from Word, to Lang) (string, error) {
	url := fmt.Sprintf("http://dict.youdao.com/search?q=%v&keyfrom=dict.index", from.W)

	content, err := gotang.HttpGetAsString(url)
	if err != nil {
		return "", errors.Wrapf(err, "请求163翻译服务出错, url: %v", url)
	}

	return t.extract(content)
}

func (t Translator163) extract(content string) (string, error) {
	start := strings.Index(content, "trans-container")
	if start <= 0 {
		return "", t.parseHtmlError("无法定位trans-container")
	}

	content = content[start:]
	start = strings.Index(content, "<li>")
	if start < 0 {
		return "", t.parseHtmlError("无法定位<li>")
	}

	content = content[start:]

	start = len("<li>")
	content = content[start:]

	end := strings.Index(content, "</li>")
	if end < 0 {
		return "", t.parseHtmlError("无法定位</li>")
	}

	return content[:end], nil
}

func (t Translator163) parseHtmlError(s string) error {
	return errors.Errorf(`解析html出错了, %v.
	 请确认是否输入了不存在的单词`)
}

type TangcloudDictRecorder struct {
	UpstreamURL string
}

func (t *TangcloudDictRecorder) Record(from Word, to Word) {
	if from.L != LANG_EN || to.L != LANG_CN {
		fmt.Printf("不支持的from %d or to %d lang", from.L, to.L)
		return
	}

	tryPostToCloud(t.UpstreamURL, from.W, to.W)
}

const MAX_TO_CHARS = 100

//TODO: 超时机制使用context.Context
func tryPostToCloud(upstreamURL, from, to string) {
	fmt.Printf("\ntry post to cloud: %s...\n", upstreamURL)
	if len(to) > MAX_TO_CHARS {
		fmt.Printf("INFO: Too large content(%v bytes), ignore post.\n", len(to))
		return
	}

	done := make(chan Result)
	gotang.Time(func() {
		timer := time.NewTimer(time.Millisecond * 2000)

		go httpPostAsString(upstreamURL /*"s"*/, postRequest{From: from, To: to}, done)

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
