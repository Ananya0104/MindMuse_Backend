package models

// ChatMessage represents a single message in a chat session
// Partition Key: userId, Sort Key: sessionId#timestamp
// This allows efficient queries for all messages by user and session, ordered by time.
type ChatMessage struct {
	UserId             string `json:"userId" dynamodbav:"userId"`           // Partition Key
	SessionId          string `json:"sessionId" dynamodbav:"sessionId"`     // Session identifier
	Timestamp          int64  `json:"timestamp" dynamodbav:"timestamp"`     // Timestamp
	SessionIdTimestamp string `json:"sessionId_timestamp" dynamodbav:"sessionId_timestamp"` // Composite sort key
	Sender             string `json:"sender" dynamodbav:"sender"`           // "user" or "ai"
	Message            string `json:"message" dynamodbav:"message"`         // Message content
} 