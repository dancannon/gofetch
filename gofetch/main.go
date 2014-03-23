package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dancannon/gofetch"
	"github.com/dancannon/gofetch/config"
	"os"
)

var (
	configFile string
	verbose    bool
	help       bool
)

func init() {
	flag.StringVar(&configFile, "config", "../config.json", "config file location")
	flag.BoolVar(&verbose, "v", false, "enable verbose logging")
	flag.BoolVar(&help, "help", false, "print help")
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s [url]:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	u := flag.Arg(0)

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	fetcher := gofetch.NewFetcher(conf)
	res, err := fetcher.Fetch(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	b, err := json.MarshalIndent(res.Content, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
