package main

import (
	"net/http"
	"rest-api-in-gin/internal/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// login logs a user in the application providing a bearer token
//
// @Summary 				Logs a user in the application providing a bearer token
// @Description 		Logs a user in the application providing a bearer token through the JWT Login Method
// @Tags 						Auth
// @Accept 					json
// @Produce 				json
// @Param 					credentials body				loginRequest	true	"Login Credentials"
// @Success 				200 					{object} 	loginResponse
// @Failure					400						{string}	string "Invalid requisition payload"
// @Failure					401						{string}	string "Invalid email or password"
// @Failure 				500 					{string}	string "Something went wrong"
// @Router /api/v1/auth/login [POST]
func (app *application) login(c *gin.Context) {
	var auth *loginRequest

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"expr":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: tokenString})
}

// registerUser registers a user in the database
//
// @Summary 				Registers a user in the database
// @Description 		Registers a user in the database
// @Tags 						Auth
// @Accept 					json
// @Produce 				json
// @Param 					credentials body				registerRequest	true	"Register credentials"
// @Success 				201 					{object} 	database.User
// @Failure					400						{string}	string "Invalid requisition payload"
// @Failure					500						{string}	string "Something went wrong"
// @Router /api/v1/auth/register [POST]
func (app *application) registerUser(c *gin.Context) {
	var register registerRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// we dont want to store the password in plain text so we gonna hash it using the bcrypt library

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	}

	register.Password = string(hashedPassword)
	user := database.User{
		Email:    register.Email,
		Password: register.Password,
		Name:     register.Name,
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
	}

	c.JSON(http.StatusCreated, user)

}
