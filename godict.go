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
	"net"
)

const (
	LangEn = iota
	LangCn
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
	Record(from Word, to Word) (ret string, err error)
}

type Translator163 struct {
}

var _ Translator = (*Translator163)(nil)

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
	return errors.Errorf(`解析html出错了:%v, 请确认是否输入了不存在的单词`, s)
}

type TangCloudDictRecorder struct {
	UpstreamURL string
}

var _ Record = (*TangCloudDictRecorder)(nil)

func (t *TangCloudDictRecorder) Record(from Word, to Word) (ret string, err error) {
	if from.L != LangEn || to.L != LangCn {
		fmt.Printf("不支持的from %d or to %d lang", from.L, to.L)
		return
	}

	return tryPostToCloud(t.UpstreamURL, from.W, to.W)
}

const MaxChars = 99

//TODO: 超时机制使用context.Context
func tryPostToCloud(upstreamURL, from, to string) (ret string, err error) {
	fmt.Printf("\ntry post to cloud: %s ...\n", upstreamURL)

	if len(to) > MaxChars {
		fmt.Printf("INFO: Too large content(%v bytes), ignore post.\n", len(to))
		return
	}

	return httpPostAsString(upstreamURL, postRequest{From: from, To: to})
}

type postRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func httpPostAsString(url string, req postRequest) (content string, err error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(req)
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport, //避免连接复用!!
	}

	resp, err := netClient.Post(url, "application/json; charset=utf-8", &buf)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(c), nil
}
