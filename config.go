package main

import "github.com/AstroX11/user-bot/types"

var AppConfig = types.Config{
	UserName: "",
	UserPN:   "",
	BotName:  "",
	Heroku: types.HerokuConfig{
		AppName: "",
		APIKey:  "",
	},
	Koyeb: types.KoyebConfig{
		AppName: "",
		APIKey:  "",
	},
	Branch: "core",
}
