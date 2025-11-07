package public

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var adminCollection *mongo.Collection

func AdminCollect() {
	adminCollection = utils.MongoClient.Database("Event_Booking").Collection("admin")
}

var (
	adminJwtKey = []byte(config.AppConfig.JWTKEY)
	adminUrl    = config.AppConfig.URL
)

func GenerateAdminToken(length int) string {
	d := make([]byte, length)
	_, _ = rand.Read(d)
	return hex.EncodeToString(d)
}

func GenerateAdminRefreshToken() string {
	return GenerateAdminToken(32)
}

// -------------------- ADMIN SIGNUP --------------------
func AdminSignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type AdminSignUpInput struct {
		AdminName string `json:"adminname" form:"adminname"`
		Email     string `json:"email" form:"email"`
		Password  string `json:"password" form:"password"`
		Phone     string `json:"phone" form:"phone"`
		Language  string `json:"language" form:"language"`
		Location  string `json:"location" form:"location"`
	}

	var input AdminSignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Request"})
		return
	}

	if input.AdminName == "" || input.Email == "" || input.Password == "" ||
		input.Phone == "" || input.Language == "" || input.Location == "" {
		c.JSON(400, gin.H{"msg": "Please fill all fields⚠️"})
		return
	}

	if !strings.Contains(input.Email, "@") || len(input.Password) < 6 || len(input.Phone) < 10 {
		c.JSON(400, gin.H{"msg": "Invalid email/password/phone"})
		return
	}

	var (
		count     int64
		countErr  error
		hashPass  []byte
		hashErr   error
		insertErr error
	)

	emailToken := GenerateAdminToken(8)
	phoneToken := GenerateAdminToken(8)

	var newAdmin models.Admin
	newAdmin.ID = primitive.NewObjectID()
	newAdmin.Role = "admin"
	newAdmin.AdminName = input.AdminName
	newAdmin.Email = input.Email
	newAdmin.Location = input.Location
	newAdmin.Language = input.Language
	newAdmin.Phone = input.Phone
	newAdmin.AdminVerified.Email = false
	newAdmin.AdminVerifyToken.Email = emailToken
	newAdmin.AdminVerifyToken.Phone = phoneToken
	newAdmin.CreatedAt = time.Now()
	newAdmin.UpdatedAt = time.Now()

	// Run count + hash concurrently
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		count, countErr = adminCollection.CountDocuments(ctx, bson.M{"email": input.Email})
	}()

	go func() {
		defer wg.Done()
		hashPass, hashErr = bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	}()

	wg.Wait()

	if countErr != nil || hashErr != nil {
		c.JSON(500, gin.H{"msg": "Internal error"})
		return
	}
	if count > 0 {
		c.JSON(400, gin.H{"msg": "Admin already exists⚠️"})
		return
	}

	newAdmin.Password = string(hashPass)
	_, insertErr = adminCollection.InsertOne(ctx, newAdmin)
	if insertErr != nil {
		c.JSON(500, gin.H{"msg": "Database error"})
		return
	}

	// Fire email in background
	go func() {
		emailData := utils.EmailData{
			From:    "Team Ivents Plannerz🎉",
			To:      input.Email,
			Subject: "Email Verification",
			Html:    fmt.Sprintf(`<a href="%s/api/public/admin/emailverify/%s">Verify email</a>`, adminUrl, emailToken),
		}
		_ = utils.SendEmail(emailData)
	}()

	c.JSON(200, gin.H{"msg": "Admin Signed Up🎉, Verify Your Email and then login✅"})
}

// -------------------- EMAIL VERIFY --------------------
func EmailVerifyAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token := c.Param("token")

	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"adminVerifyToken.emailVerifyToken": token}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Token"})
		return
	}

	if admin.AdminVerified.Email {
		c.JSON(200, gin.H{"msg": "Email already verified✅"})
		return
	}

	update := bson.M{"$set": bson.M{
		"adminverified.emailVerified":       true,
		"adminVerifyToken.emailVerifyToken": nil,
		"updated_at":                        time.Now(),
	}}

	_, err = adminCollection.UpdateByID(ctx, admin.ID, update)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Database update error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Email Verified Successfully✨"})
}

// -------------------- SIGN IN --------------------
func AdminSignIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type AdminSignInInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input AdminSignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid request"})
		return
	}

	if input.Email == "" || input.Password == "" ||
		!strings.Contains(input.Email, "@") || len(input.Password) < 6 {
		c.JSON(400, gin.H{"msg": "Invalid email/password"})
		return
	}

	var admin models.Admin
	if err := adminCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&admin); err != nil {
		c.JSON(400, gin.H{"msg": "No admin found!"})
		return
	}

	var (
		passErr       error
		tokenErr      error
		updateErr     error
		accessToken   string
		refreshToken  string
	)

	wg := sync.WaitGroup{}
	wg.Add(3)

	// Password check
	go func() {
		defer wg.Done()
		passErr = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password))
	}()

	// Access token generation
	go func() {
		defer wg.Done()
		accessToken, tokenErr = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    admin.ID.Hex(),
			"role":  admin.Role,
			"email": admin.Email,
			"exp":   time.Now().Add(5 * time.Hour).Unix(),
		}).SignedString(adminJwtKey)
	}()

	// Refresh token + DB update
	go func() {
		defer wg.Done()
		refreshToken = GenerateAdminRefreshToken()
		_, updateErr = adminCollection.UpdateByID(ctx, admin.ID, bson.M{
			"$set": bson.M{
				"refreshToken":  refreshToken,
				"refreshExpiry": time.Now().Add(7 * 24 * time.Hour),
				"updated_at":    time.Now(),
			},
		})
	}()

	wg.Wait()

	if passErr != nil {
		c.JSON(400, gin.H{"msg": "Invalid password"})
		return
	}
	if tokenErr != nil || updateErr != nil {
		c.JSON(500, gin.H{"msg": "Internal error"})
		return
	}

	c.JSON(200, gin.H{
		"msg":          "Admin logged in successfully✨",
		"token":        accessToken,
		"refreshToken": refreshToken,
	})
}

// -------------------- REFRESH TOKEN --------------------
func AdminRefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type RefreshInput struct {
		RefreshToken string `json:"refreshToken"`
	}
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		c.JSON(400, gin.H{"msg": "Invalid request"})
		return
	}

	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"refreshToken": input.RefreshToken}).Decode(&admin)
	if err != nil || admin.RefreshExpiry.IsZero() || admin.RefreshExpiry.Before(time.Now()) {
		c.JSON(401, gin.H{"msg": "Invalid or expired refresh token"})
		return
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    admin.ID.Hex(),
		"role":  admin.Role,
		"email": admin.Email,
		"exp":   time.Now().Add(5 * time.Hour).Unix(),
	}).SignedString(adminJwtKey)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Token generation failed"})
		return
	}

	c.JSON(200, gin.H{
		"msg":   "New access token generated✅",
		"token": accessToken,
	})
}

// -------------------- FORGOT PASSWORD --------------------
func AdminForgotPass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type ForgotPassInput struct {
		Email string `json:"email"`
	}

	var input ForgotPassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Request"})
		return
	}

	if input.Email == "" || !strings.Contains(input.Email, "@") {
		c.JSON(400, gin.H{"msg": "Invalid Email"})
		return
	}

	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Email not found"})
		return
	}

	tempPass := GenerateAdminToken(10)
	hashPass, hashErr := bcrypt.GenerateFromPassword([]byte(tempPass), 10)
	if hashErr != nil {
		c.JSON(500, gin.H{"msg": "Error hashing password"})
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var updateErr error
	go func() {
		defer wg.Done()
		_, updateErr = adminCollection.UpdateByID(ctx, admin.ID, bson.M{
			"$set": bson.M{
				"password":   string(hashPass),
				"updated_at": time.Now(),
			},
		})
	}()
	wg.Wait()

	if updateErr != nil {
		c.JSON(500, gin.H{"msg": "Database update failed"})
		return
	}

	go func(email, pass string) {
		emailData := utils.EmailData{
			From:    "Team Ivents Plannerz🎉",
			To:      email,
			Subject: "Password Reset Request",
			Html: fmt.Sprintf(`
				<h2>Your temporary password:</h2>
				<p style="font-size:18px; font-weight:bold;">%s</p>
				<p>Please login and change it immediately.</p>`, pass),
		}
		_ = utils.SendEmail(emailData)
	}(input.Email, tempPass)

	c.JSON(200, gin.H{"msg": "Temporary password sent to your email✅"})
}
