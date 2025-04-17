package main

import (
	"block-test/app"
	"block-test/config"
	"flag"
)

// 플래그 설정
var (
	configFlag = flag.String("environment", "./environment.toml", "environment toml file not found")
	difficulty = flag.Int("difficulty", 22, "difficulty err")
)

func main() {
	flag.Parse()
	c := config.NewConfig(*configFlag)
	app.NewApp(c)
}
