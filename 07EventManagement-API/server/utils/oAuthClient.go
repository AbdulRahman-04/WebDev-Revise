package utils

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// 👤 User Google OAuth Config
var GoogleOauthConfigUser = &oauth2.Config{
	ClientID:     config.AppConfig.OAuth.GoogleUser.ClientID,
	ClientSecret: config.AppConfig.OAuth.GoogleUser.ClientSecret,
	RedirectURL:  config.AppConfig.OAuth.GoogleUser.RedirectURL,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

// 🛡️ Admin Google OAuth Config
var GoogleOauthConfigAdmin = &oauth2.Config{
	ClientID:     config.AppConfig.OAuth.GoogleAdmin.ClientID,
	ClientSecret: config.AppConfig.OAuth.GoogleAdmin.ClientSecret,
	RedirectURL:  config.AppConfig.OAuth.GoogleAdmin.RedirectURL,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}
