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
	"golang.org/x/oauth2" // ✅ important import
)

// 🔹 Redirect USER to Google login
func GoogleLoginUser(c *gin.Context) {
	url := utils.GoogleOauthConfigUser.AuthCodeURL("user_login")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// 🔹 Redirect ADMIN to Google login
func GoogleLoginAdmin(c *gin.Context) {
	url := utils.GoogleOauthConfigAdmin.AuthCodeURL("admin_login")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// 🔹 Common callback handler (works for both user/admin)
func googleCallback(c *gin.Context, role string, oauthConfig *oauth2.Config) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code found"})
		return
	}

	// 1️⃣ Exchange code for token
	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	// 2️⃣ Get user info from Google
	client := oauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var gUser map[string]interface{}
	json.Unmarshal(data, &gUser)

	email := fmt.Sprintf("%v", gUser["email"])
	name := fmt.Sprintf("%v", gUser["name"])
	picture := fmt.Sprintf("%v", gUser["picture"])

	// 3️⃣ MongoDB connection
	userColl := utils.MongoClient.Database("Event_Booking").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing models.User
	err = userColl.FindOne(ctx, bson.M{"email": email}).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		// 🆕 Insert new user
		newUser := models.User{
			Username:  name,
			Email:     email,
			Password:  "", // Google handles authentication
			Phone:     "",
			Role:      role,
			Language:  "English",
			Location:  "Not specified",
			Createdat: time.Now(),
			Updatedat: time.Now(),
		}
		_, err := userColl.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new user"})
			return
		}
		fmt.Println("🆕 New Google", role, "added:", email)
	} else if err == nil {
		fmt.Println("✅ Existing Google", role, ":", email)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB lookup failed"})
		return
	}

	// 4️⃣ Generate JWT token
	jwtToken, err := middleware.GenerateJWT(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed"})
		return
	}

	// 5️⃣ Final success response
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Google %s login success", role),
		"email":   email,
		"name":    name,
		"avatar":  picture,
		"role":    role,
		"token":   jwtToken,
	})
}

// 🔹 User callback
func GoogleCallbackUser(c *gin.Context) {
	googleCallback(c, "user", utils.GoogleOauthConfigUser)
}

// 🔹 Admin callback
func GoogleCallbackAdmin(c *gin.Context) {
	googleCallback(c, "admin", utils.GoogleOauthConfigAdmin)
}
