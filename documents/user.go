package documents

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	UUID     string `json:"UUID"`
	LoggedAs string `json:"loggedAs"`
}

func GetUser(c *gin.Context) User {
	var user User

	claimMap, ok := c.Get("user")
	if !ok {
		return user
	}

	claim, oke := claimMap.(jwt.MapClaims)
	if !oke {
		return user
	}

	for key, val := range claim {
		switch key {
		case "UUID":
			user.UUID = val.(string)
		case "loggedAs":
			user.LoggedAs = val.(string)
		}
	}
	return user
}
