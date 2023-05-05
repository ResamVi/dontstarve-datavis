package alert

import (
	"fmt"
	"os"

	"github.com/ecnepsnai/discord"
	log "github.com/sirupsen/logrus"
)

// String sends a discord bot message to
// my notification channel for me to notice the server is down.
func String(content string) {
	log.Info(content)

	e := os.Getenv("DISCORD_WEBHOOK_URL")
	discord.WebhookURL = e
	_ = discord.Say(content)
}

func Stringf(content string, a ...any) {
	log.Info(fmt.Sprintf(content, a))

	e := os.Getenv("DISCORD_WEBHOOK_URL")
	discord.WebhookURL = e
	_ = discord.Say(fmt.Sprintf(content, a))
}

// Error sends a discord bot message to my
// notification channel for me to notice the server is down.
func Error(err error) {
	log.Errorln(err.Error())

	e := os.Getenv("DISCORD_WEBHOOK_URL")
	discord.WebhookURL = e
	_ = discord.Say(err.Error())
}
