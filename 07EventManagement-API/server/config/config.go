package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	Port    int
	DBURI   string
	URL     string
	JWTKEY  string
	Email   EmailConfig
	Phone   PhoneConfig
	Redis   RedisConfig
}

type EmailConfig struct {
	User string
	Pass string
}

type PhoneConfig struct {
	Sid   string
	Token string
	Phone string
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

var AppConfig *Config

func init() {
	// ✅ Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found")
	}

	AppConfig = &Config{
		AppName: os.Getenv("APP_NAME"),
		Port:    getEnvAsInt("PORT", 4040),
		DBURI:   os.Getenv("MONGO_URI"),
		URL:     os.Getenv("BASE_URL"),
		JWTKEY:  os.Getenv("JWT_KEY"),
		Email: EmailConfig{
			User: os.Getenv("EMAIL_USER"),
			Pass: os.Getenv("EMAIL_PASS"),
		},
		Phone: PhoneConfig{
			Sid:   os.Getenv("TWILIO_SID"),
			Token: os.Getenv("TWILIO_TOKEN"),
			Phone: os.Getenv("TWILIO_PHONE"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASS"),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

		// ✅ Add this line for debugging
	log.Println("✅ Loaded GEMINI_API_KEY:", os.Getenv("GEMINI_API_KEY"))
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
