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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed fetching user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var googleUser struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	_ = json.Unmarshal(body, &googleUser)

	// SELECT COLLECTION BASED ON ROLE
	collection := utils.MongoClient.Database("Event_Booking").Collection(role)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ==========================
	// 📌 ADMIN LOGIN FLOW
	// ==========================
	if role == "admin" {

		var existingAdmin models.Admin
		err := collection.FindOne(ctx, bson.M{"email": googleUser.Email}).Decode(&existingAdmin)

		if err == nil {
			// Already exists → JWT return
			token, _ := middleware.GenerateOAuthJWT(existingAdmin.ID.Hex(), googleUser.Email, "admin")
			c.JSON(200, gin.H{"msg": "Admin Login Success", "token": token, "role": "admin", "email": googleUser.Email})
			return
		}

		if err != mongo.ErrNoDocuments {
			c.JSON(500, gin.H{"error": "DB lookup error"})
			return
		}

		// CREATE NEW ADMIN (Lite)
		newAdmin := models.Admin{
			ID:        primitive.NewObjectID(),
			Role:      "admin",
			AdminName: googleUser.Name,
			Email:     googleUser.Email,
			Password:  "", // empty - optional login later
			Phone:     "",
			Language:  "",
			Location:  "",
			AdminVerified: struct {
				Email bool `bson:"emailVerified" json:"emailVerified"`
			}{Email: true},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = collection.InsertOne(ctx, newAdmin)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed creating admin"})
			return
		}

		token, _ := middleware.GenerateOAuthJWT(newAdmin.ID.Hex(), googleUser.Email, "admin")

		c.JSON(200, gin.H{
			"msg":   "Admin Account Created & Logged In",
			"token": token,
			"role":  "admin",
			"email": googleUser.Email,
			"name":  googleUser.Name,
		})

		return
	}

	// ==========================
	// 👤 USER LOGIN FLOW
	// ==========================
	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"email": googleUser.Email}).Decode(&existingUser)

	if err == nil {
		token, _ := middleware.GenerateOAuthJWT(existingUser.ID.Hex(), googleUser.Email, "user")
		c.JSON(200, gin.H{"msg": "User Login Success", "token": token, "role": "user"})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(500, gin.H{"error": "DB lookup error"})
		return
	}

	newUser := models.User{
		ID:         primitive.NewObjectID(),
		Role:       "user",
		Username:   googleUser.Name,
		Email:      googleUser.Email,
		Password:   "",
		Phone:      "",
		Language:   "",
		Location:   "",
		Provider:   "google",
		ProfilePic: googleUser.Picture,
		Userverified: struct {
			Email bool `bson:"emailVerified" json:"emailVerified"`
		}{Email: true},
		Createdat: time.Now(),
		Updatedat: time.Now(),
	}

	_, err = collection.InsertOne(ctx, newUser)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed creating user"})
		return
	}

	tkn, errTok := middleware.GenerateOAuthJWT(newUser.ID.Hex(), googleUser.Email, "user")
	if errTok != nil {
		c.JSON(500, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(200, gin.H{
		"msg":   "User Created & Logged In",
		"token": tkn,
		"role":  "user",
	})

}

func GoogleCallbackUser(c *gin.Context) {
	googleCallback(c, "user", utils.GoogleOauthConfigUser)
}

func GoogleCallbackAdmin(c *gin.Context) {
	googleCallback(c, "admin", utils.GoogleOauthConfigAdmin)
}
