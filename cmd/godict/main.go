package main

import (
	"fmt"
	"log"
	"os"

	"github.com/iris-contrib/color"
	"github.com/itang/godict"
)

func main() {
	word := doParseWordFromArgs()
	from := godict.Word{W: word, L: godict.LANG_EN}

	t := godict.Translator163{}

	ret, err := t.Translate(from, godict.LANG_CN)
	if err != nil {
		log.Fatalln("出错了", color.RedString(err.Error()))
	}

	fmt.Println(">", color.BlueString(ret))
}

func doParseWordFromArgs() string {
	if len(os.Args) < 2 {
		log.Fatalln(">", color.RedString("请输入要翻译的词汇"))
	}

	word := os.Args[1]
	return word
}
