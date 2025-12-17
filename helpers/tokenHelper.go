package helpers

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/yaikob/goland-jwt/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add( 24 * time.Hour).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		Uid:uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add( 7 * 24 * time.Hour).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
        return "", "", err
    }

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
        return "", "", err
    }

	return token, refreshToken, nil
}
