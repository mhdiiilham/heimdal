package config

import "github.com/spf13/viper"

type Config struct {
	RedirectURL string  `mapstructure:"redirectUrl"`
	Spotify     Spotify `mapstructure:"spotify"`
}

type Spotify struct {
	AuthorizationScopes string `mapstructure:"authorizationScopes"`
	ClientID            string `mapstructure:"clientID"`
	ClientSecret        string `mapstructure:"clientSecret"`
}

var Configuration Config

func Init() (err error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.SetConfigFile("config.yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&Configuration)
	return
}
