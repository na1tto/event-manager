package main

import (
	"rest-api-in-gin/internal/database"

	"github.com/gin-gonic/gin"
)

// util method to help extract the user from the context in the middleware

func (app *application) getUserFromContext(c *gin.Context) *database.User{
	contextUser, exist := c.Get("user")
	if !exist{
		return &database.User{}
	}
	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}

	return user
}