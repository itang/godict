package godict

import (
	"fmt"
	"strings"

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
	Record(from Lang, to Lang) error
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
