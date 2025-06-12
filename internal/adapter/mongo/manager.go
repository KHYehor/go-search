package mongo

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-search/internal/config"
	"go-search/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DBManager struct {
	client     *mongo.Client
	url        string
	db         string
	collection string
}

func NewDBManager(cfg *config.Config) *DBManager {
	return &DBManager{url: cfg.MongoURI, db: cfg.MongoDBName, collection: cfg.MongoDBCollection}
}

func (dbm *DBManager) InitConnection() error {
	client, err := mongo.Connect(options.Client().ApplyURI(dbm.url))
	if err != nil {
		return err
	}
	dbm.client = client

	return nil
}

func (dbm *DBManager) InsertItem(c *gin.Context, id string, duration int64, item any) error {
	collection := dbm.client.Database(dbm.db).Collection(dbm.collection)
	sr := entity.SearchResult{
		Id:         id,
		DurationMs: duration,
		Data:       item,
	}
	_, err := collection.InsertOne(c, sr)
	return err
}

func (dbm *DBManager) GetItem(c *gin.Context, id string) (*entity.SearchResult, error) {
	collection := dbm.client.Database(dbm.db).Collection(dbm.collection)

	sr := &entity.SearchResult{}
	err := collection.FindOne(c, bson.M{"_id": id}).Decode(sr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return sr, nil
}

func (dbm *DBManager) GetAllItems(c *gin.Context) ([]string, error) {
	collection := dbm.client.Database(dbm.db).Collection(dbm.collection)

	opts := options.Find().SetProjection(bson.M{"_id": 1})
	cursor, err := collection.Find(c, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	// Extract all at once
	var results []struct {
		ID string `bson:"_id"`
	}
	if err = cursor.All(c, &results); err != nil {
		return nil, err
	}

	// Extract IDs into a slice
	ids := make([]string, len(results))
	for i, res := range results {
		ids[i] = res.ID
	}

	return ids, nil
}
