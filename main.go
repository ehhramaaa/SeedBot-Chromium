package main

import (
	"SeedBot/core"
	"SeedBot/tools"
	"os"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func init() {
	configPath := "configs"
	proxyPath := "configs/proxy.txt"

	if !tools.CheckFileOrFolderExits(configPath) {
		os.MkdirAll(configPath, os.ModeDir)
	}

	if !tools.CheckFileOrFolderExits(proxyPath) {
		os.Create(proxyPath)
	}
}

func main() {
	defer tools.HandleRecover()

	tools.PrintLogo()

	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("configs/config.yml")
	if err != nil {
		panic(err)
	}

	core.LaunchBot()
}
