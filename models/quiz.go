package models

type QuizEntry struct {
	UserID              string            `dynamodbav:"UserID"`              // User ID from user table
	Timestamp           int64             `dynamodbav:"Timestamp"`           // Unix epoch time
	QuizQuestionnaireID string            `dynamodbav:"QuizQuestionnaireID"` // Questionnaire ID
	Answers             map[string]string `dynamodbav:"Answers"`             // question number -> user answer
}

type QuizQuestionnaire struct {
	QuizQuestionnaireID string   `dynamodbav:"QuizQuestionnaireID"` // Questionnaire ID
	Questions           []string `dynamodbav:"Questions"`           // List of questions
}
