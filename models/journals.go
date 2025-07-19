package models

// Journal represents a journal entry stored in DynamoDB
// Partition Key: UserId, Sort Key: CreatedAt
// This allows efficient queries for all journals by user and reverse chronological order.
type Journal struct {
	UserId    string `json:"userId" dynamodbav:"UserId"`       // Partition Key
	CreatedAt int64  `json:"createdAt" dynamodbav:"CreatedAt"` // Sort Key (reverse chronological)
	JournalID string `json:"journalId" dynamodbav:"JournalId"` // Unique per journal
	Date      string `json:"date" dynamodbav:"date"`
	Title     string `json:"title" dynamodbav:"title"`
	Content   string `json:"content" dynamodbav:"content"`
	UpdatedAt int64  `json:"updatedAt" dynamodbav:"updatedAt"`
}

// JournalCreateRequest represents the request body for creating a journal entry
type JournalCreateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// JournalUpdateRequest represents the request body for updating a journal entry
type JournalUpdateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// JournalResponse represents the response body for a single journal entry
// (can be extended for additional metadata if needed)
type JournalResponse struct {
	Journal Journal `json:"journal"`
	Message string  `json:"message,omitempty"`
}

// JournalListResponse represents the response body for multiple journal entries
type JournalListResponse struct {
	Journals []Journal `json:"journals"`
	Count    int       `json:"count"`
	Message  string    `json:"message,omitempty"`
}

// ErrorResponse represents a standard error response for the API
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
