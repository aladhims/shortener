package shorten

import (
	"context"
	"time"

	pb "github.com/aladhims/shortener/pkg/shorten/proto"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const databaseName = "shorten"
const collectionName = "shortURL"

type mongoRepo struct {
	c   *mongo.Client
	ctx context.Context
}

func NewMongoRepository(c *mongo.Client) Repository {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return &mongoRepo{
		c:   c,
		ctx: ctx,
	}
}

func (m *mongoRepo) GetByOrigin(origin string) *pb.ShortURL {
	var result *pb.ShortURL

	collection := m.c.Database(databaseName).Collection(collectionName)
	err := collection.FindOne(m.ctx, bson.M{"origin": origin}).Decode(result)
	if err != nil {
		return nil
	}

	return result
}

func (m *mongoRepo) GetBySlug(slug string) *pb.ShortURL {
	var result *pb.ShortURL

	collection := m.c.Database(databaseName).Collection(collectionName)
	err := collection.FindOne(m.ctx, bson.M{"slug": slug}).Decode(result)
	if err != nil {
		return nil
	}

	return result
}

func (m *mongoRepo) GetByBase(base uint64) *pb.ShortURL {
	var result *pb.ShortURL

	return result
}

func (m *mongoRepo) Create(shortURL *pb.ShortURL) error {
	collection := m.c.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(m.ctx, shortURL)
	if err != nil {
		return err
	}

	return nil

}

func (m *mongoRepo) Ping() error {
	return m.c.Ping(m.ctx, nil)
}
