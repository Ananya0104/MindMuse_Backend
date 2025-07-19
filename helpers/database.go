package helpers

import (
	"context"
	"errors"
	"fmt"
	"lambda-server/constants"
	"lambda-server/database"
	"lambda-server/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue" // V2 attribute marshalling
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types" // V2 DynamoDB types
)

var dynamoClient = database.GetInitializedClient()

/**
*   User Related DB functions
 */

// createUser creates a new user in DynamoDB
func CreateUser(user *models.User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(constants.UsersTable),
	}

	_, err = dynamoClient.PutItem(context.TODO(), input)
	if err != nil {
		print("Put Item failed with error ", err)
	}

	return err
}

// updateUser updates existing user in DynamoDB
func UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now().Unix()
	return CreateUser(user) // PutItem will overwrite
}

// getUserByID retrieves user by user ID
func GetUserByID(userID string) (*models.User, error) {

	input := &dynamodb.GetItemInput{
		TableName: aws.String(constants.UsersTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userID}, // Assuming "ID" is the primary key name
		},
	}

	result, err := dynamoClient.GetItem(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return nil, errors.New("user not found")
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	return &user, err
}

// getUserByEmail retrieves user by email using GSI
func GetUserByEmail(email string) (*models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(constants.UsersTable),
		IndexName:              aws.String("email-index"),    // Ensure this GSI exists in DynamoDB
		KeyConditionExpression: aws.String("email = :email"), // Assuming "Email" is the GSI partition key
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := dynamoClient.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("user not found")
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	return &user, err
}

// getUserByPhone retrieves user by phone number using GSI
func GetUserByPhone(phoneNumber string) (*models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(constants.UsersTable),
		IndexName:              aws.String("phoneNumber-index"),          // Ensure this GSI exists in DynamoDB
		KeyConditionExpression: aws.String("phoneNumber = :phoneNumber"), // Assuming "Phone" is the GSI partition key
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":phoneNumber": &types.AttributeValueMemberS{Value: phoneNumber},
		},
	}

	result, err := dynamoClient.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("user not found")
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	return &user, err
}

// getUserByGoogleID retrieves user by Google ID using GSI
func GetUserByGoogleID(googleID string) (*models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(constants.UsersTable),
		IndexName:              aws.String("googleId-index"),
		KeyConditionExpression: aws.String("googleId = :googleId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":googleId": &types.AttributeValueMemberS{Value: googleID},
		},
	}

	result, err := dynamoClient.Query(context.TODO(), input)
	// for k := range result.Items[0] {
	// 		fmt.Println("Key:", k)
	// }
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("user not found")
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	return &user, err
}

// FindUserByResetToken retrieves a user by password reset token
func FindUserByResetToken(token string) (*models.User, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(constants.UsersTable),
		FilterExpression: aws.String("passwordResetToken = :token"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":token": &types.AttributeValueMemberS{Value: token},
		},
	}
	result, err := dynamoClient.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, errors.New("user not found")
	}
	var user models.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	return &user, err
}

// DeleteUser deletes a user from DynamoDB by userId
func DeleteUser(userID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(constants.UsersTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userID},
		},
	}
	_, err := dynamoClient.DeleteItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete user from DynamoDB: %w", err)
	}
	return nil
}

/**
*   Survey Related DB functions
 */



// Update Survey - To update survey with user Id

/**
*   Chat Related DB functions
 */

// StoreChatMessage stores a chat message in DynamoDB
func StoreChatMessage(msg *models.ChatMessage) error {
	msg.SessionIdTimestamp = msg.SessionId + "#" + fmt.Sprintf("%d", msg.Timestamp)
	av, err := attributevalue.MarshalMap(msg)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(constants.ChatTable),
	}

	_, err = dynamoClient.PutItem(context.TODO(), input)
	if err != nil {
		print("Put Item failed with error ", err)
	}

	return err
}

// GetChatHistoryBySession retrieves all chat messages for a user and session, ordered by timestamp
func GetChatHistoryBySession(userId, sessionId string, limit int32) ([]models.ChatMessage, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(constants.ChatTable),
		KeyConditionExpression: aws.String("userId = :userId AND begins_with(sessionId_timestamp, :sessionIdPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId":          &types.AttributeValueMemberS{Value: userId},
			":sessionIdPrefix": &types.AttributeValueMemberS{Value: sessionId + "#"},
		},
		Limit:            &limit,
		ScanIndexForward: aws.Bool(true), // chronological order
	}

	result, err := dynamoClient.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	chatHistory := []models.ChatMessage{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &chatHistory)
	return chatHistory, err
}
