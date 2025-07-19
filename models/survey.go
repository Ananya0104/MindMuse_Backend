package models

type SurveyEntry struct {
	Question string   `json:"question" dynamodbav:"question"`
	Answer   string    `json:"answer" dynamodbav:"answer"`
	Elements []string `json:"elements" dynamodbav:"elements"`
}

type UpdateSurveyUserRequest struct {
	SurveyId string `json:"surveyId" binding:"required"`
	UserId   string `json:"userId" binding:"required"`
}

// Object to store response to survey
type Survey struct {
	SurveyName     string        `json:"surveyName" dynamodbav:"surveyName"`                   // Name of the survey
	SurveyId       string        `json:"surveyId" dynamodbav:"surveyId"`                       // Survey ID
	UserId         string        `json:"userId,omitempty" dynamodbav:"userId,omitempty"`       // User ID
	Timestamp      int64         `json:"timestamp" dynamodbav:"timestamp"`                     // Unix epoch time of creation of survey
	Month          string        `json:"month" dynamodbav:"month"`                             // Month of survey- Partion key - format YYYYMM - eg: 202506
	Platform       string        `json:"platform,omitempty" dynamodbav:"platform,omitempty"`   // Platform used to fill survey
	QuestionAnswer []SurveyEntry `json:"questionAnswer,omitempty" dynamodbav:"questionAnswer"` // question -> answer map
}
