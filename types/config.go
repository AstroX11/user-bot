package types

type HerokuConfig struct {
	AppName string `json:"app_name"`
	APIKey  string `json:"api_key"`
}

type KoyebConfig struct {
	AppName string `json:"app_name"`
	APIKey  string `json:"api_key"`
}

type Config struct {
	UserName string       `json:"user_name"`
	UserPN   string       `json:"user_pn"`
	BotName  string       `json:"bot_name"`
	Heroku   HerokuConfig `json:"heroku"`
	Koyeb    KoyebConfig  `json:"koyeb"`
	Branch   string       `json:"branch"`
}
