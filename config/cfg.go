package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddress    string
	SecretKey        string
	DatabaseURI      string
	DatabaseName     string
	AccessTokenTime  string
	RefreshTokenTime string
}

var Conf *Config

func init() {
	err := godotenv.Load("cfg.env")
	if err != nil {
		log.Println(err)
	}
	Conf = &Config{
		ListenAddress:    getEnv("LISTEN_ADDRESS", "localhost:8080"),
		SecretKey:        getEnv("SECRET_KEY", "jwtKEY"),
		DatabaseURI:      getEnv("DATABASE_URI", "mongodb://localhost:27017"),
		DatabaseName:     getEnv("DATABASE_NAME", "TaskJWT"),
		AccessTokenTime:  getEnv("ACCESS_TOKEN_TIME", "1m"),
		RefreshTokenTime: getEnv("REFRESH_TOKEN_TIME", "24h"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
