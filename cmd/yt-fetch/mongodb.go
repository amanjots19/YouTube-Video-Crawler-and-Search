package main

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database   = "centralDB"
	collection = "videos"
)

type mongoDB struct {
	videosColl *mongo.Collection
}

func newMongoDI(dic *diContainer) (*mongoDB, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI("mongodb+srv://aman:amanjots1@cluster0.od37c6u.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return &mongoDB{
		videosColl: client.Database(database).Collection(collection),
	}, nil
}

func newMongoDIProvider(dic *diContainer) func() (*mongoDB, error) {
	var s *mongoDB
	var mu sync.Mutex
	return func() (*mongoDB, error) {
		mu.Lock()
		defer mu.Unlock()
		var err error
		if s == nil {
			s, err = newMongoDI(dic)
		}
		return s, err
	}
}

func (m *mongoDB) getLatestDateTime() (*Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//query to find the latest video with a datetime less than the current time
	query := bson.M{"date_time": bson.M{"$lt": time.Now()}}
	opts := options.FindOne().SetSort(bson.D{{Key: "date_time", Value: -1}})

	video := new(Video)
	err := m.videosColl.FindOne(ctx, query, opts).Decode(&video)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return video, nil
		} else {
			return nil, err
		}
	}
	return video, nil
}

func (m *mongoDB) insert(v Video) error {
	_, err := m.videosColl.InsertOne(context.TODO(), v)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoDB) get(limit, offset int64) ([]Video, error) {
	findOptions := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.D{{Key: "date_time", Value: -1}})
	cur, err := m.videosColl.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}

	var videos []Video
	for cur.Next(context.TODO()) {
		var video Video
		err := cur.Decode(&video)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())
	return videos, nil
}
func (m *mongoDB) getQueryMatch(query string, limit, offset int64) ([]Video, error) {
	findOptions := options.Find().SetLimit(limit).SetSkip(offset)
	queryWords := strings.Fields(query)
	for i, word := range queryWords {
		queryWords[i] = regexp.QuoteMeta(word)
	}
	queryRegex := ".*" + strings.Join(queryWords, ".*") + ".*"

	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": queryRegex, "$options": "i"}},
			{"description": bson.M{"$regex": queryRegex, "$options": "i"}},
		},
	}
	ctx := context.TODO()
	var videos []Video
	cur, err := m.videosColl.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var video Video
		if err := cur.Decode(&video); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}
