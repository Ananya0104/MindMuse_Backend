package models

type MoodEntry struct {
	UserID              string            `dynamodbav:"UserID"`              // User ID from user table
	Timestamp           int64             `dynamodbav:"Timestamp"`           // Unix epoch time
	MoodQuestionnaireID string            `dynamodbav:"MoodQuestionnaireID"` // Questionnaire ID
	Answers             map[string]string `dynamodbav:"Answers"`             // question number -> user answer
}

type MoodQuestionnaire struct {
	MoodQuestionnaireID string   `dynamodbav:"MoodQuestionnaireID"` // Questionnaire ID
	Questions           []string `dynamodbav:"Questions"`           // List of questions
}
