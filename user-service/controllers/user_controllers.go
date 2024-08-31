package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abidaziz9876/user-service/helpers"
	"github.com/abidaziz9876/user-service/models"
	"github.com/abidaziz9876/user-service/repository"
	"github.com/abidaziz9876/user-service/response"
	"github.com/abidaziz9876/user-service/services"
	"github.com/abidaziz9876/user-service/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var User models.UserModel
		err := ctx.BindJSON(&User)
		if err != nil {
			log.Error("can not bind with json")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		err = services.CreateUser(postgres, User)
		if err != nil {
			log.Error("Error occurred " + err.Error())
			helpers.ReturnResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		helpers.ReturnResponse(ctx, http.StatusOK, "user created successfully", User.Email)
	}
}

func GetUser(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Query("id")
		if ID == "" {
			log.Error("id not found in the query")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "id not found in the query", nil)
			return
		}
		var id int64
		if i, err := strconv.ParseInt(ID, 10, 64); err == nil {
			id = i
		} else {
			log.Error("invalid id format")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "invalid id format", nil)
			return
		}

		result, err := postgres.FindUserByID(id)
		if err != nil {
			log.Error("Error occurred " + err.Error())
			helpers.ReturnResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		helpers.ReturnResponse(ctx, http.StatusOK, "user details retrieved successfully", result)

	}
}

func UpdateUser(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Query("id")
		if ID == "" {
			log.Error("id not found in the query")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "id not found in the query", nil)
			return
		}

		var user models.UserModel
		err := ctx.BindJSON(&user)
		if err != nil {
			log.Error("can not bind with json")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}

		var id int64
		if i, err := strconv.ParseInt(ID, 10, 64); err == nil {
			id = i
		} else {
			log.Error("invalid id format")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "invalid id format", nil)
			return
		}

		existingUser, err := postgres.FindUserByID(id)
		if err != nil {
			log.Error("Error occurred " + err.Error())
			helpers.ReturnResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		if existingUser == nil {
			log.Error("user not found")
			helpers.ReturnResponse(ctx, http.StatusNotFound, "user not found", nil)
			return
		}

		// Update the existing user's details
		existingUser.FirstName = user.FirstName
		existingUser.LastName = user.LastName
		existingUser.Phone = user.Phone
		existingUser.Email = user.Email
		existingUser.Password = user.Password

		result := postgres.UpdateUser(existingUser)
		if result.Error != nil {
			log.Error("Error occurred " + result.Error.Error())
			helpers.ReturnResponse(ctx, http.StatusInternalServerError, result.Error.Error(), nil)
			return
		}

		helpers.ReturnResponse(ctx, http.StatusOK, "user updated successfully", existingUser)
	}
}

func DeleteUser(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID := ctx.Query("id")
		if ID == "" {
			log.Error("id not found in the query")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "id not found in the query", nil)
			return
		}

		var id int64
		if i, err := strconv.ParseInt(ID, 10, 64); err == nil {
			id = i
		} else {
			log.Error("invalid id format")
			helpers.ReturnResponse(ctx, http.StatusBadRequest, "invalid id format", nil)
			return
		}

		existingUser, err := postgres.FindUserByID(id)
		if err != nil {
			log.Error("Error occurred " + err.Error())
			helpers.ReturnResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		if existingUser == nil {
			log.Error("user not found")
			helpers.ReturnResponse(ctx, http.StatusNotFound, "user not found", nil)
			return
		}

		result := postgres.DeleteUserByID(id)
		if result.Error != nil {
			log.Error("Error occurred " + result.Error.Error())
			helpers.ReturnResponse(ctx, http.StatusInternalServerError, result.Error.Error(), nil)
			return
		}

		helpers.ReturnResponse(ctx, http.StatusOK, "user deleted successfully", nil)
	}
}

var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func SignUp(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.UserModel
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		existingUser, err := postgres.FindUserByEmailOrPhone(user.Email, user.Phone)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if existingUser != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("user already exists")})
			return
		}

		password := HashPassword(string(user.Password))
		user.Password = password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		token, refreshtoken, _ := tokens.TokenGenerator(user.Email, user.FirstName, user.LastName)
		user.Token = token
		user.Refresh_Token = refreshtoken
		result := postgres.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, response.ApiResponse{
			Status:  200,
			Data:    user,
			Message: "Successfully SignUp",
		})
	}
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

func LogIn(postgres repository.PostgesDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.UserModel
		var founduser models.UserModel
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		query := `select * from ticketbooking.users u where u.email =?`
		res := postgres.RawQueryWithFind(&founduser, query, user.Email)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(user.Password, founduser.Password)

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := tokens.TokenGenerator(founduser.Email, founduser.FirstName, founduser.LastName)

		tokens.UpdateAllTokens(token, refreshToken, founduser.ID)
		c.JSON(http.StatusOK, response.ApiResponse{
			Status:  200,
			Data:    founduser,
			Message: "Successfully Logged In",
		})
	}
}


func Check() gin.HandlerFunc  {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,"Hello world")
	}
}