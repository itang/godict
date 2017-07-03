package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/itang/godict"
	"github.com/pkg/errors"
)

type config struct {
	UpstreamURL string `toml:"upstream_url"`
}

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
	upstreamURL, err := getUpstreamURL()
	if err != nil {
		upstreamURL = "http://www.godocking.com/api/dict/log"
	}

	var record godict.Record = &godict.TangcloudDictRecorder{UpstreamURL: upstreamURL}
	record.Record(godict.Word{W: word, L: godict.LANG_EN}, godict.Word{W: ret, L: godict.LANG_CN})
}

func getUpstreamURL() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", nil
	}
	tomlData, err := ioutil.ReadFile(usr.HomeDir + "/.rdict/config.toml") // just pass the file name
	if err != nil {
		return "", nil
	}
	var conf config
	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		return "", err
	}
	return conf.UpstreamURL, nil
}

func parseWordFromArgs() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("请输入要翻译的词汇")
	}

	word := os.Args[1]
	return word, nil
}
