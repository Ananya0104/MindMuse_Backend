package constants

const (
	// DynamoDB Table Names
	UsersTable         string = "mindmuse_users"
	JournalsTable      string = "mindmuse_journal"
	JournalQueryLimit  int    = 20
	MindMuseScoreTable string = "mindmuse_score"

	// Add chat table name
	ChatTable string = "mindmuse_chat"

	// DynamoDB Key Names for Journals
	DynamoDbKeyUserId    string = "UserId"
	DynamoDbKeyJournalId string = "JournalId"
)

// To identify if project is running locally or on cloud
const (
	IsRunningLocally string = "is_running_locally"
	IsEnvSet         string = "is_env_set"
	PathToEnv        string = "env.yaml"
)

// Query parameter and context key names
const (
	QueryParamJournalId string = "journalId"
	QueryParamUserId    string = "userId"
	ContextKeyUserId    string = "userId"
)
