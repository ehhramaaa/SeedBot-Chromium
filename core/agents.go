package core

import (
	"SeedBot/tools"

	"github.com/mileusna/useragent"
)

func randomUserAgent() (string, string) {
	userAgents, err := tools.ReadFileTxt("./configs/useragent.txt")
	if err != nil {
		tools.Logger("error", err.Error())
	}

	userAgent := userAgents[tools.RandomNumber(0, len(userAgents))]

	os := useragent.Parse(userAgent).OS

	return userAgent, os
}
