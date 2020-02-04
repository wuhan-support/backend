package main

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Shimoauth struct {
		ClientId     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		Scope        string `yaml:"scope"`
	} `yaml:"shimoauth"`
	DBInfo struct {
		DBName string `yaml:"dbname"`
		Addr   string `yaml:"addr"`
		User   string `yaml:"user"`
		Pwd    string `yaml:"pwd"`
	} `yaml:"dbinfo"`
	UploadPath string `yaml:"uploadPath"`
	Telegram struct {
		BotToken string `yaml:"botToken"`
		ChatID int64 `yaml:"chatId"`
	} `yaml:"telegram"`
}
