package repository

import (
	"TASKONE/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository interface {
	Insert(log model.UserLog) error
	GetByUserID(userID int) ([]model.UserLog, error)
}

type logRepo struct {
	Collection mongo.Collection
}

func NewLogRepository(db *mongo.Database) LogRepository {
	return &logRepo{Collection: *db.Collection("user_logs")}
}

func (r *logRepo) Insert(log model.UserLog) error {
	log.Timestamp = time.Now()
	_, err := r.Collection.InsertOne(context.Background(), log)
	return err
}

func (r *logRepo) GetByUserID(userID int) ([]model.UserLog, error) {
	var logs []model.UserLog
	cursor, err := r.Collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var l model.UserLog
		cursor.Decode(&l)
		logs = append(logs, l)
	}
	return logs, nil
}
