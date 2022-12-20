package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/lucsky/cuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("email").NotEmpty().Unique(),
		field.String("Password").NotEmpty().DefaultFunc(cuid.New),
		field.Bool("IsAdmin").Default(false),
		field.String("name").Default("unknown").Optional(),
		field.String("first_name").Default("unknown").Optional(),
		field.String("nickname").Default("unknown").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("Services", Service.Type).Ref("users"),
	}
}
