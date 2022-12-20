package models

import (
	"api-gateway/database"
	"api-gateway/ent"
	"api-gateway/ent/service"

)

type ServicesModel struct {
	Server   ServerModel `json:"server"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

type AboutModel struct {
	Client ClientModel `json:"client"`
	Server ServerModel `json:"server"`
}

type ClientModel struct {
	Host string `json:"host"`
}

type ServerModel struct {
	CurrentTime int64          `json:"current_time"`
	Services    []ServiceModel `json:"services"`
}

type ServiceModel struct {
	Name      string             `json:"name"`
	Actions   []InteractionModel `json:"actions"`
	Reactions []InteractionModel `json:"reactions"`
}

type InteractionModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserServicesRequest struct {
	Services []string `json:"services"`
}

func GetUserServices(user *ent.User) ([]*ent.Service, error) {
	services, err := user.QueryServices().
		All(database.Db.Ctx)
	return services, err
}

func CreateService(services ServiceModel) (*ent.Service, error) {

	service, err := database.Db.Def.Service.
		Create().
		SetName(services.Name).
		Save(database.Db.Ctx)
	return service, err
}

func GetServices() ([]*ent.Service, error) {
	services, err := database.Db.Def.Service.
		Query().
		All(database.Db.Ctx)
	return services, err
}

func GetServicesByName(name string) (*ent.Service, error) {
	services, err := database.Db.Def.Service.
		Query().
		Where(service.Name(name)).
		Only(database.Db.Ctx)
	return services, err
}

func DeleteServices() (int, error) {
	services, err := database.Db.Def.Service.
		Delete().Exec(database.Db.Ctx)
	return services, err
}
