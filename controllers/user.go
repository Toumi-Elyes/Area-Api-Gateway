// Package controllers provide controllers for the backend routes
package controllers

import (
	"api-gateway/logger"
	"api-gateway/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var APISECRET = []byte(os.Getenv("API_SECRET_KEY"))

// GenerateToken Create json web token using the email of the user
func GenerateToken(user models.UserModel) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = user.Email
	claims["user_id"] = user.Id
	claims["user_password"] = user.Password
	claims["exp"] = time.Now().Add(time.Minute * 1440).Unix() //Token expires after 1 day
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(APISECRET)
}

func ExtractToken(c *gin.Context) (string, string, error) {
	userToken := c.GetHeader("Authorization")
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return APISECRET, nil
	})
	if err != nil {
		return "", "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return fmt.Sprintf("%v", claims["user_id"]), fmt.Sprintf("%v", claims["user_email"]), err
	}
	return "", "", nil
}

// Register Route: POST /createUser
// Create a new user and send the corresponding response.
// :Body: Json with the email of the user and his password ({"email": "user@example.com", "password": "example"}).
// On success: return 200 with the user in json format.
// if the Json is not well formatted, returns 400 (bad request).
// if the user already exists, or any error in the database, returns 406 (status not acceptable).
func Register(c *gin.Context) {
	userData := models.UserModel{}
	err := c.ShouldBindBodyWith(&userData, binding.JSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	isEmail := models.CheckEmail(userData)
	if !isEmail {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}
	isHashed := models.HashPassword(&userData)

	if !isHashed {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}

	user, err := models.CreateUser(userData)

	if err != nil {
		logger.Log.Errorf("Failed to create the user in the database.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusConflict)
	}

	userData.Id = user.ID
	token, err := GenerateToken(userData)

	if err != nil {
		logger.Log.Errorf("Failed to generate user's token.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err != nil {
		logger.Log.Errorf("Failed to create new User.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusConflict)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

func Session(c *gin.Context) {
	c.AbortWithStatus(http.StatusAccepted)
}

// DeleteUser Route: POST "/deleteUser"
// Delete a user and send the corresponding response. On success: return 200 with the string "User has been deleted successfully".
// :Body: Json with the id of the user ({"id": "ex-am-pl-e"}).
// On success: return 200 with the following message: "User has been deleted successfully".
// if the Json is not well formatted, returns 400 (bad request).
// if the user has already been deleted, or any error in the database, returns 406 (status not acceptable).
func DeleteUser(c *gin.Context) {
	userData := models.UserModel{}
	err := c.ShouldBindBodyWith(&userData, binding.JSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	err = models.DeleteUser(userData)

	if err != nil {
		logger.Log.Errorf("Failed to delete the User.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	} else {
		c.JSON(http.StatusOK, fmt.Sprint("User ", userData.Id, " has been deleted successfully"))
	}
}

// UpdateUser Route: POST "/updateUser"
// Update the user information and send the corresponding response
// :Body: Json with the id, the email and the password of the user ({"id": "ex-am-pl-e", "email": "user@example.com", "password": "example"}).
// On success: return 200 with the user in json format.
// if the Json is not well formatted, returns 400 (bad request).
// if the user does not exist, or any error in the database, returns 406 (status not acceptable).
func UpdateUser(c *gin.Context) {
	userData := models.UserModel{}
	err := c.ShouldBindBodyWith(&userData, binding.JSON)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	isEmail := models.CheckEmail(userData)
	if !isEmail {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}
	isHashed := models.HashPassword(&userData)

	if !isHashed {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}
	tokenId, _, err := ExtractToken(c)

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	tmpId, err := uuid.Parse(tokenId)

	if err != nil {
		c.String(http.StatusInternalServerError, "Cannot convert string to uuid")
		return
	}
	userData.Id = tmpId
	user, err := models.UpdateUser(userData)

	userData.Id = user.ID
	userData.Email = user.Email
	userData.Password = user.Password
	token, err := GenerateToken(userData)
	if err != nil {
		logger.Log.Errorf("Failed to update the User.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

// ReadUser Route: POST "/readUser"
// Read one user and send the corresponding response
// :Body: Json with the id of the user ({"id": "ex-am-pl-e"}).
// On success: return 200 with the user information in json format.
// if the Json is not well formatted, returns 400 (bad request).
// if the user does not exist, or any error in the database, returns 406 (status not acceptable).
func ReadUser(c *gin.Context) {
	tokenId, _, err := ExtractToken(c)

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	tmpId, err := uuid.Parse(tokenId)

	fmt.Println(tokenId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Cannot convert string to uuid")
		return
	}

	logger.Log.Errorf("%v\n", tmpId)
	fmt.Println(tmpId)
	userData := models.UserModel{}
	userData.Id = tmpId

	logger.Log.Errorf("%v\n", userData)
	user, err := models.ReadOneUser(userData)

	if err != nil {
		logger.Log.Errorf("Failed to read the User.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// ReadUsers Route: POST "/readUsers"
// Read all the users in the database and send the corresponding response
// :Body: Admin access required.
// On success: return 200 with the user information as a list of json elements.
// if the user does not exist, or any error in the database, returns 406 (status not acceptable).
func ReadUsers(c *gin.Context) {
	users, err := models.ReadUsers()

	if err != nil {
		logger.Log.Errorf("Failed to read the Users.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	} else {
		c.JSON(http.StatusOK, users)
	}
}

// Login Route: POST "/login"
// check if the user exist and send the jwt on success
// :Body: Admin access required.
// On success: return 200 with the users' information as a list of json elements.
// if the user does not exist, or any error in the database, returns 406 (status not acceptable).
func Login(c *gin.Context) {
	user := models.UserModel{}
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	users, err := models.ReadUsers()

	if err != nil {
		logger.Log.Errorf("Failed to read the User.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	}

	isHashed := models.HashPassword(&user)

	if !isHashed {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}

	for _, u := range users {
		if u.Email == user.Email {
			if u.Password != user.Password {
				c.String(http.StatusNotAcceptable, "Wrong password")
				return
			}
			user.Id = u.ID
			token, err := GenerateToken(user)
			if err != nil {
				c.String(http.StatusInternalServerError, "Cannot generate JWT Token")
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
			return
		}
	}
	c.String(http.StatusBadRequest, "Bad Request")
}

// loginAdmin
func LoginAdmin(c *gin.Context) {
	user := models.UserModel{}
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	users, err := models.ReadUsers()

	if err != nil {
		logger.Log.Errorf("Failed to read the User.\nError message: [%v]\n", err)
		c.AbortWithStatus(http.StatusNotAcceptable)
	}

	isHashed := models.HashPassword(&user)

	if !isHashed {
		c.AbortWithStatus(http.StatusNotAcceptable)
	}

	for _, u := range users {
		if u.Email == user.Email {
			if u.IsAdmin == false {
				c.String(http.StatusNotAcceptable, "You are not an admin")
				return
			}
			if u.Password != user.Password {
				c.String(http.StatusNotAcceptable, "Wrong password")
				return
			}
			user.Id = u.ID
			token, err := GenerateToken(user)
			if err != nil {
				c.String(http.StatusInternalServerError, "Cannot generate JWT Token")
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
			return
		}
	}
	c.String(http.StatusBadRequest, "Bad Request")
}
