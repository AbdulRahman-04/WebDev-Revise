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

// Redirect USER to Google login
func GoogleLoginUser(c *gin.Context) {
	url := utils.GoogleOauthConfigUser.AuthCodeURL("user_login")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Redirect ADMIN to Google login
func GoogleLoginAdmin(c *gin.Context) {
	url := utils.GoogleOauthConfigAdmin.AuthCodeURL("admin_login")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Common callback (used by both)
func googleCallback(c *gin.Context, role string, oauthConfig *oauth2.Config) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code found"})
		return
	}

	// Step 1: Exchange code for token
	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	// Step 2: Fetch user info from Google
	client := oauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var gUser map[string]interface{}
	json.Unmarshal(body, &gUser)

	email := fmt.Sprintf("%v", gUser["email"])
	name := fmt.Sprintf("%v", gUser["name"])
	picture := fmt.Sprintf("%v", gUser["picture"])

	// Step 3: Choose collection based on role
	var collectionName string
	if role == "admin" {
		collectionName = "admin"
	} else {
		collectionName = "user"
	}

	userColl := utils.MongoClient.Database("Event_Booking").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 4: Check if already exists
	var existing models.User
	err = userColl.FindOne(ctx, bson.M{"email": email}).Decode(&existing)

	if err == nil {
		// ⚠️ Different response per role
		var errMsg string
		if role == "admin" {
			errMsg = "Admin already exists. Please login normally."
		} else {
			errMsg = "User already exists. Please login normally."
		}

		c.JSON(http.StatusConflict, gin.H{
			"error":   errMsg,
			"message": fmt.Sprintf("Same %s email already registered.", role),
		})
		fmt.Println("⚠️ Duplicate Google login attempt for:", email, "| Role:", role)
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB lookup failed"})
		return
	}

	// Step 5: Insert new Google user/admin
	newUser := models.User{
		Username:   name,
		Email:      email,
		Password:   "",
		Phone:      "",
		Role:       role,
		Language:   "English",
		Location:   "Not specified",
		Provider:   "google",
		ProfilePic: picture,
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

	fmt.Println("🆕 New Google", role, "added:", email)

	// Step 6: Generate JWT
	jwtToken, err := middleware.GenerateJWT(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed"})
		return
	}

	// Step 7: Send success response
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Google %s login success", role),
		"name":    name,
		"email":   email,
		"avatar":  picture,
		"role":    role,
		"token":   jwtToken,
	})
}

// User callback
func GoogleCallbackUser(c *gin.Context) {
	googleCallback(c, "user", utils.GoogleOauthConfigUser)
}

// Admin callback
func GoogleCallbackAdmin(c *gin.Context) {
	googleCallback(c, "admin", utils.GoogleOauthConfigAdmin)
}
