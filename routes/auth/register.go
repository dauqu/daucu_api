package auth

import (
	"context"
	"daucu/config"
	"daucu/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UsersCollection *mongo.Collection = config.GetCollection(config.DB, "users")

func Register(c *gin.Context) {

	var users models.Users
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Check all fields
	if users.Fullname == "" || users.Email == "" || users.Username == "" || users.Password == "" {
		c.JSON(400, gin.H{"message": "All fields are required"})
		return
	}

	//Check if username exists
	var user models.Users
	err := UsersCollection.FindOne(ctx, bson.M{"username": users.Username}).Decode(&user)
	if err != mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"message": "Username already exists"})
		return
	}

	//Check if email exists
	err = UsersCollection.FindOne(ctx, bson.M{"email": users.Email}).Decode(&user)
	if err != mongo.ErrNoDocuments {
		c.JSON(400, gin.H{"message": "Email already exists"})
		return
	}

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), 8)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Insert user
	_, err = UsersCollection.InsertOne(ctx, bson.M{
		"fullname": users.Fullname,
		"email":    users.Email,
		"username": users.Username,
		"password": hashedPassword,
		"role":     "admin",
		"license":  "free",
	})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created",
	})
}
