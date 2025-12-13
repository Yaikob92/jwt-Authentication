package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"git.hub.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"github.com/yaikob/goland-jwt/database"
	"github.com/yaikob/goland-jwt/helpers"
	"github.com/yaikob/goland-jwt/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword(){
}

func VerifyPassword(){}

func Signup(){
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var use model.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if ValidationErr != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationError.Error()})
			return
		}
		count,err:= userCollection.CountDocmuments(ctx,bson.M{"email":user.Email})
		defer cancel()
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H("error":"error accored while checking for the email"))
			return
		}
		count,err := userCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		defer cancel()
		if err != nil{
			log.panic(err)
			c.JSON(http.StatusInternalServerError,gin.H("error":"error accored while checking for the phone number"))
			return
		}
		if count>0{
			c.JSON(http.StatusInternalServerError,gin.H("error":"this email or phone number already exists"))
			return
		}

	}
}

func Login(){}

func GetUsers(){}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Param("user_id")
		if err := helpers.MatchUserTypeToUid(c,userId); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var user models.User
		err  := userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}
		c.JSON(http.StatusOK,user)


	}

}
