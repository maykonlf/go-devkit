package uuid

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type sample struct {
	ID UUID
}

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

func TestMongoDBUUIDMarshaller(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017/"))
	assert.Nil(t, err)

	c := client.Database("testdb").Collection("uuid_test")
	id := New()

	_, err = c.InsertOne(context.Background(), &sample{ID: id})
	assert.Nil(t, err)

	var fetchedSample sample
	err = c.FindOne(context.Background(), &bson.D{{"id", id}}).Decode(&fetchedSample)
	assert.NotNil(t, err)
	assert.Equal(t, id.String(), fetchedSample.ID.String())
}
