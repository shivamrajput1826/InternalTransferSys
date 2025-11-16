package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() {
	profile := viper.GetString("PROFILE")
	if profile == "" {
		profile = "dev"
	}

	envFile := "./config/dev.env"
	godotenv.Load(envFile)

	viper.AutomaticEnv()

	viper.SetConfigName(profile)
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../config/")

	_ = viper.ReadInConfig()
}

func GetConfigValues(key string) string {
	return viper.GetString(key)
}
func GetPort() string {
	return viper.GetString("port")
}
