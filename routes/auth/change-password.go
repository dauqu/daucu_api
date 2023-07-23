package auth

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(c *gin.Context) {

	//Get token from header
	token := c.GetHeader("dauqu-auth-token")

	//Check if token is empty
	if token == "" {
		c.JSON(400, gin.H{"message": "Token is empty"})
		return
	}

	type Body struct {
		Username    string `json:"username"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Find user
	user, err := UsersCollection.FindOne(ctx, bson.M{"username": body.Username}).DecodeBytes()
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Check if user is empty
	if user == nil {
		c.JSON(400, gin.H{"message": "User not found"})
		return
	}

	//Get hashed password
	hashedPassword := user.Lookup("password").StringValue()

	//Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body.OldPassword)); err != nil {
		c.JSON(400, gin.H{"message": "Incorrect password"})
		return
	}

	//Hash new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 8)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Update user password
	_, err = UsersCollection.UpdateOne(ctx, bson.M{"username": body.Username}, bson.M{"$set": bson.M{"password": newHashedPassword}})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Password changed successfully"})
}
