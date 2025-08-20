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
		field.Int64("id").
			Unique(),
		field.Int64("uid").
			Unique().
			Positive(),
		field.String("provider_type"). // 改为枚举类型
			MaxLen(20),
		field.String("provider_id").
			MaxLen(255), // 增加长度以匹配SQL定义
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the AuthProvider.
func (AuthProvider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("auth_providers").
			Field("uid").
			Required().
			Unique(),
	}
}
