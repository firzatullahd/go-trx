package main

import (
	"flag"
	"fmt"
	"go-trx/config"
	"go-trx/logger"
	"go-trx/server"
)

func main() {
	var configName string
	flag.StringVar(&configName, "config_name", "config", "A config name that used by server")
	flag.Parse()

	conf := config.Load(configName)

	fmt.Printf("%+v", conf)

	logger.Init(conf)
	server.Start(conf)
}
