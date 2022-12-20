// Package models Model package store all the models and call the database when needed
package models

import (
	"api-gateway/database"
	"api-gateway/ent"
	"api-gateway/ent/user"
	"crypto/sha256"
	"encoding/base64"
	"net/mail"

	"github.com/google/uuid"
)

type UserModel struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Firstname string    `json:"firstname"`
	Nickname  string    `json:"nickname"`
}

func CheckEmail(userData UserModel) bool {
	_, err := mail.ParseAddress(userData.Email)
	if err != nil {
		return false
	}
	return true
}

func HashPassword(userData *UserModel) bool {
	hash := sha256.New()
	if userData.Password == "" {
		return false
	} else {
		hash.Write([]byte(userData.Password))
		userData.Password = base64.URLEncoding.EncodeToString(hash.Sum(nil))
		return true
	}
}

// CreateUser Create a new user in the database with the given email and password.
// Parameter, email and password of the user to create.
// On success, return the new User with nil as error.
// On failure, return nil for a user and the error.
func CreateUser(user UserModel) (*ent.User, error) {
	newUser, err := database.Db.Def.User.
		Create().
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetName(user.Name).
		SetFirstName(user.Firstname).
		SetNickname(user.Nickname).
		Save(database.Db.Ctx)
	return newUser, err
}

// DeleteUser Delete a user from the database.
// Parameter, id of the user to delete.
// On success, return nil.
// On failure, return the error.
func DeleteUser(data UserModel) error {
	err := database.Db.Def.User.
		DeleteOneID(data.Id).
		Exec(database.Db.Ctx)
	return err
}

// UpdateUser Create a new user in the database with the given email and password.
// Parameter, information of the user to update.
// On success, return the updated User with nil as error.
// On failure, return nil for a user and the error.
func UpdateUser(user UserModel) (*ent.User, error) {
	updatedUser, err := database.Db.Def.User.
		UpdateOneID(user.Id).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetName(user.Name).
		SetFirstName(user.Firstname).
		SetNickname(user.Nickname).
		Save(database.Db.Ctx)
	return updatedUser, err
}

// ReadOneUser Retrieve the information about one user of the database.
// Parameter, id of the user to retrieve.
// On success, return the User with nil as error.
// On failure, return nil for a user and the error.
func ReadOneUser(data UserModel) (*ent.User, error) {
	u, err := database.Db.Def.User.
		Query().
		Where(
			user.ID(data.Id)).
		Only(database.Db.Ctx)
	return u, err
}

// ReadUsers Retrieve the information of all the database's users.
// Parameter, no params.
// On success, return the users with nil as error.
// On failure, return nil for a user and the error.
func ReadUsers() ([]*ent.User, error) {
	users, err := database.Db.Def.User.
		Query().
		All(database.Db.Ctx)
	return users, err
}

func AddServiceToUser(user UserModel, service *ent.Service) error {
	_, err := database.Db.Def.User.
		UpdateOneID(user.Id).
		AddServices(service).
		Save(database.Db.Ctx)
	return err
}

func RemoveServiceFromUser(user UserModel, service *ent.Service) error {
	_, err := database.Db.Def.User.
		UpdateOneID(user.Id).
		RemoveServices(service).
		Save(database.Db.Ctx)
	return err
}
