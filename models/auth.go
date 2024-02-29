package models

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	Cash      int                `json:"cash,omitempty" bson:"cash,omitempty"`
}

func (u *User) HashPassword() {
	result, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(result)
}

func (u *User) ComparePasswords(canditatePassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(canditatePassword, []byte(u.Password))
	return err == nil
}

func (u *User) CreateJWT() (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	encriptionKey := os.Getenv("JWT_ENCRIPTION_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  u.ID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString([]byte(encriptionKey))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
