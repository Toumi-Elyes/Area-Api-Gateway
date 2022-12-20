package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type Area struct {
	ent.Schema
}

// Fields of the Area.
func (Area) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("area_id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("area_name").NotEmpty(),
		field.String("action_reaction").NotEmpty(),
	}
}

// Edges of the Area.
func (Area) Edges() []ent.Edge {
	return nil
}
