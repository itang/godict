package godict

import (
	"fmt"

	"errors"
	"github.com/itang/gotang"
	"strings"
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
	Record(from Lang, to Lang) error
}

type Translator163 struct {
}

func (t Translator163) Translate(from Word, to Lang) (string, error) {
	url := fmt.Sprintf("http://dict.youdao.com/search?q=%v&keyfrom=dict.index", from.W)

	content, err := gotang.HttpGetAsString(url)
	if err != nil {
		return "", errors.New("请求163翻译服务出错:" + err.Error())
	}

	return t.extract(content)
}

func (t Translator163) extract(content string) (string, error) {
	start := strings.Index(content, "trans-container")
	if start <= 0 {
		return "", errors.New("error1")
	}

	content = content[start:]
	start = strings.Index(content, "<li>")
	if start < 0 {
		return "", errors.New("error2")
	}

	content = content[start:]

	start = len("<li>")
	content = content[start:]

	end := strings.Index(content, "</li>")
	if end < 0 {
		return "", errors.New("error3")
	}

	return content[:end], nil
}
