package models

type MindMuseScore struct {
	UserId    string  `json:"userId" dynamodbav:"userId"`
	Score     float64 `json:"score" dynamodbav:"score"`
	Timestamp int64   `json:"timestamp" dynamodbav:"timestamp"`
} 