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
	OAuth   OAuthConfig
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

// ✅ Separate Google OAuth for user & admin
type OAuthConfig struct {
	GoogleUser  GoogleOAuth
	GoogleAdmin GoogleOAuth
	Github      GithubOAuth
}

type GoogleOAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GithubOAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var AppConfig *Config

func init() {
	// Load .env
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

		OAuth: OAuthConfig{
			GoogleUser: GoogleOAuth{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID_USER"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET_USER"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL_USER"),
			},
			GoogleAdmin: GoogleOAuth{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID_ADMIN"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET_ADMIN"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL_ADMIN"),
			},
			Github: GithubOAuth{
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			},
		},
	}

	log.Println("✅ Config Loaded | Port:", AppConfig.Port)
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
