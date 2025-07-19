package models

type SOS struct {
	UserID    string `dynamodbav:"UserID"`    // User ID from user table
	Timestamp int64  `dynamodbav:"Timestamp"` // Unix epoch time
	Status    string `dynamodbav:"Status"`
}
