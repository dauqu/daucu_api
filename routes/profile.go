package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	//"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v4"
	//"golang.org/x/crypto/bcrypt"
)

func Profile(c *gin.Context) {
	//Get authorization header
	token := c.Request.Header.Get("Authorization")

	//Parse token
	parse_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	//Check if token is valid
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	//Check if token is valid
	if !parse_token.Valid {
		c.JSON(400, gin.H{"message": "Invalid token"})
		return
	}

	//Get claims
	claims := parse_token.Claims.(jwt.MapClaims)

	//Get username from claims
	username := claims["username"].(string)

	//Create context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := UsersCollection.Find(ctx, bson.M{"username": username})
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var users []bson.M

	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, users)
}
