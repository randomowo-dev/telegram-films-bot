package database

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	TelegramID int64              `bson:"telegram_id,omitempty"`
	Username   string             `bson:"username,omitempty"`
	LastAuth   time.Time          `bson:"last_auth,omitempty"`
	Created    time.Time          `bson:"created,omitempty"`
	Updated    time.Time          `bson:"updated"`
	Role       Role               `bson:"role,omitempty"`
}

type Role uint8

func (r *Role) String() string {
	switch *r {
	case ModeratorRole:
		return "moderator"
	case AdminRole:
		return "admin"
	default:
		return "user"
	}
}

func (r *Role) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, r.String())), nil
}

func (r *Role) UnmarshalJSON(bytes []byte) error {
	*r = ParseRole(string(bytes)[1 : len(bytes)-1])
	return nil
}

func (r *Role) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(r.String())
}

func (r *Role) UnmarshalBSONValue(t bsontype.Type, bytes []byte) error {
	if t != bson.TypeString {
		return fmt.Errorf("wrong bson type for Role type")
	}

	var res string
	err := bson.UnmarshalValue(t, bytes, &res)
	if err != nil {
		return err
	}

	*r = ParseRole(res)

	return nil
}

func ParseRole(role string) Role {
	switch role {
	case "user":
		return UserRole
	case "moderator":
		return ModeratorRole
	case "admin":
		return AdminRole
	default:
		return UserRole
	}
}

const (
	NoRole Role = iota
	UserRole
	ModeratorRole
	AdminRole
)

var AvailableRoles = []Role{UserRole, ModeratorRole, AdminRole}
