package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
)

// User holds the schema definition for the User entity.
type Service struct {
	ent.Schema
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
	}
}

func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}
