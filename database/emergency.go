package database

import (
	"context"
	"errors"
	"fmt"
	"lambda-server/constants"
	"lambda-server/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue" // V2 attribute marshalling
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types" // V2 DynamoDB types
)

// SetEmergencyContacts updates the EmergencyContacts field for a user in the Users table
func SetEmergencyContacts(userID string, contacts [3]models.Emergency) error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(constants.UsersTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userID},
		},
	}
	result, err := GetInitializedClient().GetItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to get user from DynamoDB: %w", err)
	}
	if result.Item == nil {
		return errors.New("user not found")
	}
	var user models.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %w", err)
	}
	user.EmergencyContacts = contacts
	user.UpdatedAt = time.Now().Unix()
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}
	putInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(constants.UsersTable),
	}
	_, err = GetInitializedClient().PutItem(context.TODO(), putInput)
	if err != nil {
		return fmt.Errorf("failed to update user in DynamoDB: %w", err)
	}
	return nil
}

// GetEmergencyContacts retrieves the EmergencyContacts field for a user from the Users table
func GetEmergencyContacts(userID string) ([3]models.Emergency, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(constants.UsersTable),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: userID}, // Assuming "ID" is the primary key name
		},
	}
	result, err := GetInitializedClient().GetItem(context.TODO(), input)
	if err != nil {
		return [3]models.Emergency{}, fmt.Errorf("failed to get user from DynamoDB: %w", err)
	}
	if result.Item == nil {
		return [3]models.Emergency{}, errors.New("user not found")
	}
	var user models.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return [3]models.Emergency{}, fmt.Errorf("failed to unmarshal user: %w", err)
	}
	return user.EmergencyContacts, nil
}
