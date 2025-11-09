package utils

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// 👤 User GitHub OAuth Config
var GithubOauthConfigUser = &oauth2.Config{
	ClientID:     config.AppConfig.OAuth.GithubUser.ClientID,
	ClientSecret: config.AppConfig.OAuth.GithubUser.ClientSecret,
	RedirectURL:  config.AppConfig.OAuth.GithubUser.RedirectURL,
	Scopes:       []string{"read:user", "user:email"},
	Endpoint:     github.Endpoint,
}

// 🛡️ Admin GitHub OAuth Config
var GithubOauthConfigAdmin = &oauth2.Config{
	ClientID:     config.AppConfig.OAuth.GithubAdmin.ClientID,
	ClientSecret: config.AppConfig.OAuth.GithubAdmin.ClientSecret,
	RedirectURL:  config.AppConfig.OAuth.GithubAdmin.RedirectURL,
	Scopes:       []string{"read:user", "user:email"},
	Endpoint:     github.Endpoint,
}
