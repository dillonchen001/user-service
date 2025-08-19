package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AuthProvider holds the schema definition for the AuthProvider entity.
type AuthProvider struct {
	ent.Schema
}

// Fields of the AuthProvider.
func (AuthProvider) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider_type").
			MaxLen(20),
		field.String("provider_id").
			MaxLen(100),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the AuthProvider.
func (AuthProvider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("auth_providers").
			Unique(),
	}
}
