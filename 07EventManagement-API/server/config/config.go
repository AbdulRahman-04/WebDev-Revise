package config

type Config struct {
	AppName string
	Port int
	DBURL string
	URL string
	JWTKEY string
	Email EmailConfig 
	Phone PhoneConfig
	Redis RedisConfig
}

type EmailConfig struct {
	USER string
	PASS string
}

type PhoneConfig struct {
	SID string
	TOKEN string
	PHONE string
}

type RedisConfig struct {
	Host string
	Password string
	DB int
}

var AppConfig = &Config{
	AppName: "BookMyEvent.com",
	Port: 6565,
	DBURL: "mongodb+srv://abdrahman:abdrahman@rahmann18.hy9zl.mongodb.net/BookMyEvent.com",
	URL: "http://localhost:6565",
	JWTKEY: "RAHMAN123",
	Email: EmailConfig{
		USER: "abdulrahman.81869@gmail.com",
		PASS: "ttkv mljj ukcp ijcl",
	},
	Phone: PhoneConfig{
		SID: "",
		TOKEN: "",
		PHONE: "",
	},
	Redis: RedisConfig{
		Host: "localhost:6379",
		Password: "",
		DB:  0,
	},
}