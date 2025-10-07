package model

import "time"

type UserLog struct {
	UserID    int       `bson:"user_id" json:"user_id"`
	Action    string    `bson:"action" json:"action"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}
