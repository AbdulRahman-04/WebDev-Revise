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

var userCollection *mongo.Collection

func UserCollect() {
	userCollection = utils.MongoClient.Database("Event_Booking").Collection("user")
}

var userJwtKey = []byte(config.AppConfig.JWTKEY)
var userUrl = config.AppConfig.URL

func GenerateUserToken(length int) string {
	d := make([]byte, length)
	_, _ = rand.Read(d)
	return hex.EncodeToString(d)
}

func GenerateRefreshToken() string {
	return GenerateUserToken(32)
}

func UserSignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type UserSignUp struct {
		UserName string `json:"name" form:"name"`
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
		Phone    string `json:"phone" form:"phone"`
		Language string `json:"language" form:"language"`
		Location string `json:"location" form:"location"`
	}

	var inputUser UserSignUp
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Request"})
		return
	}

	if inputUser.UserName == "" || inputUser.Email == "" || inputUser.Password == "" ||
		inputUser.Phone == "" || inputUser.Language == "" || inputUser.Location == "" {
		c.JSON(400, gin.H{"msg": "Please fill all fields⚠️"})
		return
	}
	if !strings.Contains(inputUser.Email, "@") || len(inputUser.Password) < 6 || len(inputUser.Phone) < 10 {
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

	emailToken := GenerateUserToken(8)
	phoneToken := GenerateUserToken(8)

	var newUser models.User
	newUser.ID = primitive.NewObjectID()
	newUser.Role = "user"
	newUser.Username = inputUser.UserName
	newUser.Email = inputUser.Email
	newUser.Location = inputUser.Location
	newUser.Language = inputUser.Language
	newUser.Phone = inputUser.Phone
	newUser.Userverified.Email = false
	newUser.Userverifytoken.Email = emailToken
	newUser.Userverifytoken.Phone = phoneToken
	newUser.Createdat = time.Now()
	newUser.Updatedat = time.Now()

	// Run count + hash concurrently
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		count, countErr = userCollection.CountDocuments(ctx, bson.M{"email": inputUser.Email})
	}()

	go func() {
		defer wg.Done()
		hashPass, hashErr = bcrypt.GenerateFromPassword([]byte(inputUser.Password), 10)
	}()

	wg.Wait()

	if countErr != nil || hashErr != nil {
		c.JSON(500, gin.H{"msg": "Internal error"})
		return
	}
	if count > 0 {
		c.JSON(400, gin.H{"msg": "User already exists⚠️"})
		return
	}

	newUser.Password = string(hashPass)
	_, insertErr = userCollection.InsertOne(ctx, newUser)
	if insertErr != nil {
		c.JSON(500, gin.H{"msg": "Database error"})
		return
	}

	// Fire email in background (no wait)
	go func() {
		emailData := utils.EmailData{
			From:    "Team Ivents Plannerz🎉",
			To:      inputUser.Email,
			Subject: "Email Verification",
			Html:    fmt.Sprintf(`<a href="%s/api/public/user/emailverify/%s">Verify email</a>`, userUrl, emailToken),
		}
		_ = utils.SendEmail(emailData)
	}()

	c.JSON(200, gin.H{"msg": "User Signed Up🎉, Verify Your Email and then login✅"})
}

// -------------------- EMAIL VERIFY --------------------
func EmailVerifyUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token := c.Param("token")

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"userverifytoken.emailVerifyToken": token}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Token"})
		return
	}

	if user.Userverified.Email {
		c.JSON(200, gin.H{"msg": "Email Verified already, u can login now!"})
		return
	}

	update := bson.M{"$set": bson.M{
		"userverified.emailVerified":       true,
		"userverifytoken.emailVerifyToken": nil,
		"updated_at":                       time.Now(),
	}}

	_, err = userCollection.UpdateByID(ctx, user.ID, update)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "email Verified✨🙌"})
}

//                     SIGN IN API 
func UserSignIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type UserSignIn struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}

	var inputUser UserSignIn
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid request"})
		return
	}

	if inputUser.Email == "" || inputUser.Password == "" ||
		!strings.Contains(inputUser.Email, "@") || len(inputUser.Password) < 6 {
		c.JSON(400, gin.H{"msg": "Invalid email/password"})
		return
	}

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": inputUser.Email}).Decode(&user); err != nil {
		c.JSON(400, gin.H{"msg": "No email found!"})
		return
	}

	var (
		passErr    error
		tokenErr   error
		updateErr  error
		accessToken string
		refreshToken string
	)

	wg := sync.WaitGroup{}
	wg.Add(3)

	// Password check
	go func() {
		defer wg.Done()
		passErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputUser.Password))
	}()

	// Access token generation
	go func() {
		defer wg.Done()
		accessToken, tokenErr = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    user.ID.Hex(),
			"role":  user.Role,
			"email": user.Email,
			"exp":   time.Now().Add(5 * time.Hour).Unix(),
		}).SignedString(userJwtKey)
	}()

	// Refresh token + DB update
	go func() {
		defer wg.Done()
		refreshToken = GenerateRefreshToken()
		_, updateErr = userCollection.UpdateByID(ctx, user.ID, bson.M{
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
		"msg":          "Logged in successfully!✨",
		"token":        accessToken,
		"refreshToken": refreshToken,
	})
}

// -------------------- REFRESH ACCESS TOKEN --------------------
func RefreshToken(c *gin.Context) {
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

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"refreshToken": input.RefreshToken}).Decode(&user)
	if err != nil || user.RefreshExpiry.IsZero() || user.RefreshExpiry.Before(time.Now()) {
		c.JSON(401, gin.H{"msg": "Invalid or expired refresh token"})
		return
	}

	// Generate new access token
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID.Hex(),
		"role":  user.Role,
		"email": user.Email,
		"exp":   time.Now().Add(5 * time.Hour).Unix(),
	}).SignedString(userJwtKey)
	if err != nil {
		c.JSON(400, gin.H{"msg": "token generation failed"})
		return
	}

	c.JSON(200, gin.H{
		"msg":   "New access token generated",
		"token": accessToken,
	})
}

func UserChangePass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type UserChangePass struct {
		Email       string `json:"email" form:"email"`
		Oldpassword string `json:"oldpassword" form:"oldpassword"`
		Newpassword string `json:"newpassword" form:"newpassword"`
	}

	var inputUser UserChangePass
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid request"})
		return
	}
	if inputUser.Email == "" || inputUser.Oldpassword == "" || inputUser.Newpassword == "" ||
		!strings.Contains(inputUser.Email, "@") || len(inputUser.Newpassword) < 6 {
		c.JSON(400, gin.H{"msg": "Invalid input"})
		return
	}

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": inputUser.Email}).Decode(&user); err != nil {
		c.JSON(400, gin.H{"msg": "No email found"})
		return
	}

	var (
		passErr    error
		hashPass   []byte
		hashErr    error
		updateErr  error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		passErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputUser.Oldpassword))
	}()

	go func() {
		defer wg.Done()
		hashPass, hashErr = bcrypt.GenerateFromPassword([]byte(inputUser.Newpassword), 10)
	}()

	wg.Wait()
	if passErr != nil {
		c.JSON(400, gin.H{"msg": "Invalid old password"})
		return
	}
	if hashErr != nil {
		c.JSON(500, gin.H{"msg": "Hashing failed"})
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		update := bson.M{"$set": bson.M{
			"password":   string(hashPass),
			"updated_at": time.Now(),
		}}
		_, updateErr = userCollection.UpdateByID(ctx, user.ID, update)
	}()
	wg.Wait()

	if updateErr != nil {
		c.JSON(400, gin.H{"msg": "Database error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Password Changed Successfully!✅"})
}


func UserForgotPass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type UserForgotPass struct {
		Email string `json:"email" form:"email"`
	}

	var inputUser UserForgotPass
	if err := c.ShouldBindJSON(&inputUser); err != nil || !strings.Contains(inputUser.Email, "@") || inputUser.Email == "" {
		c.JSON(400, gin.H{"msg": "Invalid Email"})
		return
	}

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": inputUser.Email}).Decode(&user); err != nil {
		c.JSON(400, gin.H{"msg": "No email found"})
		return
	}

	tempPass := GenerateUserToken(8)

	var (
		hashNewPass []byte
		hashErr     error
		updateErr   error
	)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		hashNewPass, hashErr = bcrypt.GenerateFromPassword([]byte(tempPass), 10)
	}()

	wg.Wait()
	if hashErr != nil {
		c.JSON(500, gin.H{"msg": "Hashing failed"})
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		update := bson.M{"$set": bson.M{
			"password":   string(hashNewPass),
			"updated_at": time.Now(),
		}}
		_, updateErr = userCollection.UpdateByID(ctx, user.ID, update)
	}()
	wg.Wait()

	if updateErr != nil {
		c.JSON(400, gin.H{"msg": "Database error"})
		return
	}

	// 🚀 Fire email in background (no wait)
	go func() {
		emailData := utils.EmailData{
			From:    "Team Ivents Plannerz🎉",
			To:      inputUser.Email,
			Subject: "Reset Password Request",
			Html:    fmt.Sprintf(`<h2>Your Temporary password is <strong>%s</strong></h2>`, tempPass),
		}
		_ = utils.SendEmail(emailData)
	}()

	c.JSON(200, gin.H{"msg": "Temporary password sent to your email✅✨"})
}