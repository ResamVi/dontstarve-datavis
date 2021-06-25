package alert

import "github.com/ecnepsnai/discord"

func Msg(content string) {
	discord.WebhookURL = ""
	discord.Say(content)
}
