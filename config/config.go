package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() {
	// Load environment variables first
	profile := viper.GetString("PROFILE")
	if profile == "" {
		profile = "dev"
	}

	// Try to load .env file
	envFile := "./config/dev.env"
	godotenv.Load(envFile)

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Try to load YAML config if it exists
	viper.SetConfigName(profile)
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../config/")

	// Ignore error if config file doesn't exist
	_ = viper.ReadInConfig()
}

func GetConfigValues(key string) string {
	return viper.GetString(key)
}
func GetPort() string {
	return viper.GetString("port")
}
