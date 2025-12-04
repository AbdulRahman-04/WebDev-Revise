package auth

import (
	"context"
	"encoding/json"
	// "fmt"
	"io"
	"net/http"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

func GoogleLoginUser(c *gin.Context) {
	url := utils.GoogleOauthConfigUser.AuthCodeURL("user_login", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleLoginAdmin(c *gin.Context) {
	url := utils.GoogleOauthConfigAdmin.AuthCodeURL("admin_login", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func googleCallback(c *gin.Context, role string, oauthConfig *oauth2.Config) {
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

	client := oauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var googleUser struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	json.Unmarshal(body, &googleUser)

	userColl := utils.MongoClient.Database("Event_Booking").Collection(role)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// CHECK IF USER EXISTS
	var existing models.User
	err = userColl.FindOne(ctx, bson.M{"email": googleUser.Email}).Decode(&existing)

	// 🔹 If user exists → generate JWT & return
	if err == nil {
		jwtToken, _ := middleware.GenerateOAuthJWT(existing.ID.Hex(), googleUser.Email, role)

		c.JSON(http.StatusOK, gin.H{
			"message": "Login success",
			"name":    googleUser.Name,
			"email":   googleUser.Email,
			"avatar":  googleUser.Picture,
			"role":    role,
			"token":   jwtToken,
		})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB lookup error"})
		return
	}

	// INSERT NEW USER
	newUser := models.User{
		ID:         primitive.NewObjectID(),
		Username:   googleUser.Name,
		Email:      googleUser.Email,
		Password:   "",
		Phone:      "",
		Role:       role,
		Provider:   "google",
		ProfilePic: googleUser.Picture,
		Userverified: struct {
			Email bool `bson:"emailVerified" json:"emailVerified"`
		}{Email: true},
		Createdat: time.Now(),
		Updatedat: time.Now(),
	}

	_, err = userColl.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// GENERATE JWT
	jwtToken, _ := middleware.GenerateOAuthJWT(newUser.ID.Hex(), googleUser.Email, role)

	c.JSON(http.StatusOK, gin.H{
		"message": "Google " + role + " login success",
		"name":    googleUser.Name,
		"email":   googleUser.Email,
		"avatar":  googleUser.Picture,
		"role":    role,
		"token":   jwtToken,
	})
}

func GoogleCallbackUser(c *gin.Context) {
	googleCallback(c, "user", utils.GoogleOauthConfigUser)
}

func GoogleCallbackAdmin(c *gin.Context) {
	googleCallback(c, "admin", utils.GoogleOauthConfigAdmin)
}
