package config

import (
	"os"

	"github.com/AstroX11/user-bot/types"
	"github.com/joho/godotenv"
)

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

func init() {
	_ = godotenv.Load()

	if AppConfig.UserName == "" {
		AppConfig.UserName = os.Getenv("USER_NAME")
	}
	if AppConfig.UserPN == "" {
		AppConfig.UserPN = os.Getenv("USER_PN")
	}
	if AppConfig.BotName == "" {
		AppConfig.BotName = os.Getenv("BOT_NAME")
	}
	if AppConfig.Heroku.AppName == "" {
		AppConfig.Heroku.AppName = os.Getenv("HEROKU_APP_NAME")
	}
	if AppConfig.Heroku.APIKey == "" {
		AppConfig.Heroku.APIKey = os.Getenv("HEROKU_API_KEY")
	}
	if AppConfig.Koyeb.AppName == "" {
		AppConfig.Koyeb.AppName = os.Getenv("KOYEB_APP_NAME")
	}
	if AppConfig.Koyeb.APIKey == "" {
		AppConfig.Koyeb.APIKey = os.Getenv("KOYEB_API_KEY")
	}
	if AppConfig.Branch == "" {
		AppConfig.Branch = os.Getenv("BRANCH")
		if AppConfig.Branch == "" {
			AppConfig.Branch = "core"
		}
	}
}
