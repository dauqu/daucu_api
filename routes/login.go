package routes

import (
	"context"
	"time"

	"daucu/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	//JWT
	"github.com/golang-jwt/jwt/v4"
)

func Login(c *gin.Context) {

	type Login struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var login Login

	//Bind JSON
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)

	//Check if user exists
	user, err := UsersCollection.FindOne(ctx, bson.M{"username": login.Username}).DecodeBytes()
	if err == mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"message": "User not found"})
		return
	}

	//Decode user
	var u models.Users
	err = bson.Unmarshal(user, &u)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Get hashed password
	hashedPassword := u.Password

	//Check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(login.Password))
	if err != nil {
		c.JSON(400, gin.H{"message": "Incorrect password"})
		return
	}

	//Generate token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = login.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Set token in cookies same site none
	// c.SetSameSite(http.SameSiteNoneMode)
	// c.SetCookie("token", t, 86400, "/", "", true, false)

	//Return token
	c.JSON(200, gin.H{"message": "Login successful", "token": t})
}
