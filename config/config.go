package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("yaml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../config/")
	viper.AutomaticEnv()

	profile := viper.GetString("PROFILE")
	if profile == "" {
		profile = "dev"
		godotenv.Load("./config/.env.dev")
	}

	viper.SetConfigName(profile)

}

func GetConfigValues(key string) string {
	return viper.GetString(key)
}
func GetPort() string {
	return viper.GetString("port")
}
