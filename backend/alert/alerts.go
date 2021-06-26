package alert

import "github.com/ecnepsnai/discord"

func Msg(content string) {
	discord.WebhookURL = ""
	discord.Say(content)
}

func Panic(content string) {
	discord.WebhookURL = ""
	discord.Say(content)
	panic(content)
}
