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
	SID   string
	TOKEN  string
	PHONE string
}

type RedisConfig struct {
	Host string
	Password string
	DB int
}

var AppConfig = &Config{
	AppName: "BookMyEVENT",
	Port: 3030,
	DBURL: "",
	URL: "http://localhost:3030",
	JWTKEY: "RAHMAN123",
	Email: EmailConfig{
		USER:  "",
		PASS: "",
	},
	Phone: PhoneConfig{
		SID: "",
		TOKEN: "",
		PHONE: "",
	},
	Redis:  RedisConfig{
		Host: "localhots:6379",
		Password: "",
		DB: 0,
	},
}