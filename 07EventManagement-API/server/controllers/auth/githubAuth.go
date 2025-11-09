package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

// Redirect USER to GitHub login
func GithubLoginUser(c *gin.Context) {
	url := utils.GithubOauthConfigUser.AuthCodeURL("user_login")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Redirect ADMIN to GitHub login
func GithubLoginAdmin(c *gin.Context) {
	url := utils.GithubOauthConfigAdmin.AuthCodeURL("admin_login")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Common callback handler
func githubCallback(c *gin.Context, role string, oauthConfig *oauth2.Config) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code found"})
		return
	}

	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	// Fetch user info from GitHub
	client := oauthConfig.Client(c, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ghUser map[string]interface{}
	json.Unmarshal(body, &ghUser)

	email := fmt.Sprintf("%v", ghUser["email"])
	name := fmt.Sprintf("%v", ghUser["login"]) // GitHub username
	avatar := fmt.Sprintf("%v", ghUser["avatar_url"])

	// Pick correct collection
	var collectionName string
	if role == "admin" {
		collectionName = "admin"
	} else {
		collectionName = "user"
	}

	userColl := utils.MongoClient.Database("Event_Booking").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing models.User
	err = userColl.FindOne(ctx, bson.M{"email": email}).Decode(&existing)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   fmt.Sprintf("%s already exists. Please login normally.", role),
			"message": "Same email already registered.",
		})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB lookup failed"})
		return
	}

	newUser := models.User{
		Username:   name,
		Email:      email,
		Password:   "",
		Phone:      "",
		Role:       role,
		Language:   "English",
		Location:   "Not specified",
		Provider:   "github",
		ProfilePic: avatar,
		Userverified: struct {
			Email bool `bson:"emailVerified" json:"emailVerified"`
		}{Email: true},
		Createdat: time.Now(),
		Updatedat: time.Now(),
	}

	_, err = userColl.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new user"})
		return
	}

	jwtToken, err := middleware.GenerateJWT(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("GitHub %s login success", role),
		"name":    name,
		"email":   email,
		"avatar":  avatar,
		"role":    role,
		"token":   jwtToken,
	})
}

// User callback
func GithubCallbackUser(c *gin.Context) {
	githubCallback(c, "user", utils.GithubOauthConfigUser)
}

// Admin callback
func GithubCallbackAdmin(c *gin.Context) {
	githubCallback(c, "admin", utils.GithubOauthConfigAdmin)
}
