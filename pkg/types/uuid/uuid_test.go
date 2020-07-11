package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func TestUUIDMarshaling(t *testing.T) {
	uuid := New()

	bsonType, bytes, err := uuid.MarshalBSONValue()
	assert.Nil(t, err)

	unmarshalError := uuid.UnmarshalBSONValue(bsonType, bytes)
	assert.Nil(t, unmarshalError)
}

func TestUUID_UnmarshalBSONValue(t *testing.T) {
	uuid := New()

	_, bytes, err := uuid.MarshalBSONValue()
	assert.Nil(t, err)

	err = uuid.UnmarshalBSONValue(bsontype.Array, bytes)
	assert.Error(t, err, "invalid format on unmarshall bson value")

	err = uuid.UnmarshalBSONValue(bsontype.Binary, []byte(`1234`))
	assert.Error(t, err, "not enough bytes to unmarshal bson value")
}
