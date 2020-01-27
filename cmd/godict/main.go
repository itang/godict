package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/itang/godict"
	"github.com/itang/gotang"
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

	ret, err := godict.Translator163{}.Translate(godict.Word{W: word, L: godict.LangEn}, godict.LangCn)
	if err != nil {
		fmt.Printf("ERROR: %s\n", color.RedString(err.Error()))
		return
	}
	fmt.Printf("\t->: %s\n", color.BlueString(ret))

	upstreamURL, err := getUpstreamURL()
	if err != nil {
		//upstreamURL = "http://www.godocking.com/api/dict/log"
		fmt.Printf("read post url config from %v error: %v\n", configFilePath(), err)
		return
	}

	if upstreamURL == "" {
		fmt.Printf("\n\t%s\n", color.RedString("read config url config from %v, got empty, ignore to post\n", configFilePath()))
		return
	}

	var record godict.Record = &godict.TangCloudDictRecorder{UpstreamURL: upstreamURL}
	gotang.Time(func() {
		ret, err := record.Record(godict.Word{W: word, L: godict.LangEn}, godict.Word{W: ret, L: godict.LangCn})
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
		} else {
			fmt.Printf("\t|>: %v\n", ret)
		}
	})
}

func getUpstreamURL() (string, error) {
	tomlData, err := ioutil.ReadFile(configFilePath()) // just pass the file name
	if err != nil {
		return "", err
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

func configFilePath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir + "/.rdict/config.toml"
}
