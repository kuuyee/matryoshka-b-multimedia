package main

import (
	"flag"

	"github.com/kuuyee/matryoshka-b-multimedia/conf"
	"github.com/kuuyee/matryoshka-b-multimedia/runner"
)

func main() {
	confPath := flag.String("c", "multimedia.yml", "config file")
	flag.Parse()

	conf.Get(*confPath)

	runner.Run()
}
