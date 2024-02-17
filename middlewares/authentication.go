package middlewares

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// https://medium.com/@smooth-55/json-web-token-jwt-authentication-in-golang-using-jwt-go-245fd18e14af
var encriptionKey string = os.Getenv("JWT_ENCRIPTION_KEY")

func authentication(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("token")
	if tokenString == "" {
		//return unautherized
	}
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(encriptionKey), nil
	})
	if err != nil {
		// return nil, err
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if ok && token.Valid {
		// check if the token expires or not
		if float64(time.Now().Unix()) > float64(claims.ExpiresAt) {
			m.logger.Zap.Error("token expired [ParseWithClaims]")
			err := errors.BadRequest.New("Token already expired")
			err = errors.SetCustomMessage(err, "Token expired")
			return false, err
		}
	} else {
		err := errors.BadRequest.New("Invalid token")
		err = errors.SetCustomMessage(err, "Invalid token")
		// return false, err
	}

	// Get user from claims and set
	user, err := m.userService.GetOneUser(claims.Id)
	if err != nil {
		m.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		m.logger.Zap.Error("error finding user record")
		return false, err
	}
	// Can set anything in the request context and passes the request to the next handler.
	c.Set("user_id", user.ID)
	c.Set("username", user.username)
	return true, nil
}
