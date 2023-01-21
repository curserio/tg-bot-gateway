package tgbot

type Config struct {
	Verbose bool     `yaml:"verbose"`
	Token   string   `yaml:"token"`
	Admins  []string `yaml:"admins"`
}
