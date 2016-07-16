package main

import (
	"fmt"
	"log"
	"os"

	"github.com/itang/godict"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln(">请输入要翻译的词汇")
	}

	word := os.Args[1]
	from := godict.Word{W: word, L: godict.LANG_EN}

	t := godict.Translator163{}

	ret, err := t.Translate(from, godict.LANG_CN)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(">", ret)
}
