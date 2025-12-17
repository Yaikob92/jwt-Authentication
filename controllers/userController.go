package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yaikob/goland-jwt/database"
	"github.com/yaikob/goland-jwt/helpers"
	"github.com/yaikob/goland-jwt/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bcrypt.GenerateFromPassword([]byte(password),14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "email or password is incorrect"
		check = false
	}
	return check, msg

}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error accored while checking for the email"})
			return
		}

		phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error accored while checking for the phone number"})
			return
		}
		

		if emailCount > 0 || phoneCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email or phone number already exists"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, err := helpers.GenerateAllTokens(
			*user.Email,
			*user.First_name,
			*user.Last_name,
			*user.User_type,
			user.User_id,
		)

		user.Token = &token
		user.Refresh_token = &refreshToken
		if err != nil {
			log.Println("Eror generating token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to genereate token"})
			return
		}

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}
		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "email or password is incorrect"})
			return

		}
		passwordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordValid {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": msg})
			return
		}

		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"user not found"})

		}

		token,refreshToken,_ := helpers.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.First_name,
			*foundUser.Last_name,
			*foundUser.User_type,
			foundUser.User_id,
		)
		helpers.UpdateTokens(token,refreshToken,foundUser.User_id)
		err = userCollection.FindOne(ctx,bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}
		c.JSON(http.StatusOK,foundUser)
	}
}

func GetUsers() {}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)

	}

}
