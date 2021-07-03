package alert

import (
	log "github.com/sirupsen/logrus"

	"github.com/ecnepsnai/discord"
)

func Msg(content string) {
	discord.WebhookURL = ""
	discord.Say(content)
	log.Println(content)
}

func Panic(content string) {
	discord.WebhookURL = ""
	discord.Say(content)
	panic(content)
}
