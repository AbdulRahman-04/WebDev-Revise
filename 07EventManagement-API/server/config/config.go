package config

type  Config struct {
	AppName string
	Port int
	DBURL string
	URL string
	JWTKEY string
	Email EmailConfig
	Phone PhoneConfig
	Redis RedisConfig
}

type EmailConfig struct{
	USER string
	PASS string
}
 
type PhoneConfig struct{
	SID string
    TOKEN string
	PHONE string
}

type RedisConfig struct{
	Host string
	Password string
	DB int
}

var AppConfig = &Config{
	AppName: "BookMyEvents.com",
	Port: 9090,
	DBURL: "mongodb+srv://abdrahman:abdrahman@rahmann18.hy9zl.mongodb.net/BookMyEvents.com",
	URL: "http://localhost:9099",
	JWTKEY: "RAHMAN123",
	Email: EmailConfig{
		USER: "abdulrahman.81869@gmail.com",
		PASS: "gdwi eqgw bisl iihw",
	},
	Phone: PhoneConfig{
		SID: "",
		TOKEN: "",
		PHONE: "",
	},
	Redis: RedisConfig{
		Host: "localhost:6379",
		Password: "",
		DB: 0,
	},
}