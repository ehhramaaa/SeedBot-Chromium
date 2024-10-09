package core

import (
	"SeedBot/tools"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gookit/config/v2"
)

func LaunchBot() {
	defer tools.HandleRecover()
	sessionsPath := "sessions"
	proxyPath := "configs/proxy.txt"
	maxThread := config.Int("MAX_THREAD")
	isUseProxy := config.Bool("USE_PROXY")

	if !tools.CheckFileOrFolderExits(sessionsPath) {
		os.MkdirAll(sessionsPath, os.ModeDir)
	}

	sessionList, err := tools.ReadFileInDir(sessionsPath)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed To Read File Directory: %v", err))
	}

	if len(sessionList) <= 0 {
		tools.Logger("error", "No Session Found")
	}

	var wg sync.WaitGroup
	var semaphore chan struct{}
	var proxyList []string

	tools.Logger("info", fmt.Sprintf("%v Session Detected", len(sessionList)))

	if isUseProxy {
		proxyList, err = tools.ReadFileTxt(proxyPath)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("Proxy Data Not Found: %s", err))
		}

		tools.Logger("info", fmt.Sprintf("%v Proxy Detected", len(proxyList)))
	}

	if maxThread > len(sessionList) {
		semaphore = make(chan struct{}, len(sessionList))
	} else {
		semaphore = make(chan struct{}, maxThread)
	}

	for {
		totalPointsChan := make(chan float64, len(sessionList))

		for index, session := range sessionList {
			wg.Add(1)
			account := &Account{
				Phone: strings.TrimSuffix(session.Name(), ".json"),
			}

			go account.worker(&wg, &semaphore, &totalPointsChan, index, session, proxyList)
		}

		go func() {
			wg.Wait()
			close(totalPointsChan)
		}()

		var totalPoints float64

		for points := range totalPointsChan {
			totalPoints += points
		}

		tools.Logger("success", fmt.Sprintf("Total Points All Account: %.9f", (totalPoints/1e9)))

		randomSleep := tools.RandomNumber(config.Int("RANDOM_SLEEP.MIN"), config.Int("RANDOM_SLEEP.MAX"))

		tools.Logger("info", fmt.Sprintf("Launch Bot Finished | Sleep %vs Before Next Lap...", randomSleep))

		time.Sleep(time.Duration(randomSleep) * time.Second)
	}
}
