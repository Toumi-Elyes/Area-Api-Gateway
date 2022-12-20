package models

import (
	"api-gateway/database"
	"api-gateway/ent"
	"api-gateway/ent/area"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type DataResponse struct {
	Success bool   `json:"success"`
	Payload string `json:"payload"`
}

type AreaRequestModel struct {
	Payload AreaModel `json:"payload"`
}

type AreaToDelete struct {
	ServiceName string `json:"serviceName"`
	Type        string `json:"type"`
}

type AreaDeleteRequestDispatcher struct {
	AreaId   uuid.UUID      `json:"areaId"`
	AreaList []AreaToDelete `json:"areaList"`
}

type AreaModel struct {
	AreaId         uuid.UUID             `json:"area_id"`
	AreaName       string                `json:"area_name"`
	UserId         uuid.UUID             `json:"user_id"`
	ActionReaction []ActionReactionModel `json:"action_reaction"`
}

type ActionReactionModel struct {
	Type    string         `json:"type"`
	Service string         `json:"service"`
	Name    string         `json:"name"`
	Order   int            `json:"order"`
	Payload map[string]any `json:"payload"`
}

type AuthentificatorDataModel struct {
	AccessToken string `json:"access_token"`
	Service     string `json:"service"`
}

type ArOrReaField struct {
	Name  string
	Value any
}

type ArOrRea struct {
	Type    string
	Service string
	Name    string
	Order   int
	Icon    string
	Fields  []ArOrReaField
}

type UserArea struct {
	AreaName       string // name donné par le user
	AreaId         uuid.UUID
	ActionReaction []ArOrRea
}

type UserAreasProps struct {
	Areas []UserArea
}

type DeleteAreaRequest struct {
	AreaId uuid.UUID `json:"areaId"`
}

type IconResponse struct {
	Icon string `json:"icon"`
}

func CreateArea(area AreaRequestModel) (*ent.Area, error) {
	actions_reactions, _ := json.Marshal(area.Payload.ActionReaction)
	newArea, err := database.Db.Def.Area.
		Create().
		SetAreaName(area.Payload.AreaName).
		SetUserID(area.Payload.UserId).
		SetActionReaction(string(actions_reactions)).
		Save(database.Db.Ctx)
	return newArea, err
}

func GetAreas() ([]*ent.Area, error) {
	areas, err := database.Db.Def.Area.
		Query().
		All(database.Db.Ctx)
	return areas, err
}

func CheckAreas(userId uuid.UUID, areas []*ent.Area) bool {
	for i := 0; i < len(areas); i += 1 {
		if areas[i].UserID == userId {
			return true
		}
	}
	return false
}

func GetUserAreas(data UserModel) ([]*ent.Area, error) {
	areas, err := database.Db.Def.Area.
		Query().
		Where(area.UserID(data.Id)).
		All(database.Db.Ctx)
	return areas, err
}

func GetAreasById(id uuid.UUID) (*ent.Area, error) {
	areas, err := database.Db.Def.Area.
		Query().
		Where(area.AreaID(id)).
		First(database.Db.Ctx)
	return areas, err
}

func DeleteArea(id uuid.UUID) (int, error) {
	nbRow, err := database.Db.Def.Area.
		Delete().
		Where(area.AreaID(id)).
		Exec(database.Db.Ctx)
	return nbRow, err
}
