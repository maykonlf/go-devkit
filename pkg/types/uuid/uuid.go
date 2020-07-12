package uuid

import (
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type UUID struct {
	uuid.UUID
}

func New() UUID {
	return UUID{uuid.New()}
}

func (u UUID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.Binary, bsoncore.AppendBinary(nil, 4, u.UUID[:]), nil
}

func (u *UUID) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	if t != bsontype.Binary {
		return errors.New("invalid format on unmarshall bson value")
	}

	_, data, _, ok := bsoncore.ReadBinary(raw)
	if !ok {
		return errors.New("not enough bytes to unmarshal bson value")
	}

	copy(u.UUID[:], data)
	return nil
}

func MustParse(s string) UUID {
	return UUID{uuid.MustParse(s)}
}

func Parse(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	return UUID{id}, err
}

func FromBytes(b []byte) (UUID, error) {
	id, err := uuid.FromBytes(b)
	return UUID{id}, err
}
