package database

import (
	"context"
	"fmt"
	"lambda-server/constants"
	"lambda-server/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CreateJournalEntry creates a new journal entry in DynamoDB
func CreateJournalEntry(ctx context.Context, entry models.Journal) error {
	// Convert the journal entry to DynamoDB attribute values
	item, err := attributevalue.MarshalMap(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal journal entry: %w", err)
	}

	// Put the item in the table
	_, err = GetInitializedClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(constants.JournalsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

// GetJournalByID retrieves a specific journal entry by userId and journalId using the GSI
func GetJournalByID(ctx context.Context, userId string, journalId string) (*models.Journal, error) {
	// Query the GSI to get the item with userId and journalId
	result, err := GetInitializedClient().Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(constants.JournalsTable),
		IndexName:              aws.String("UserId-JournalId-index"),
		KeyConditionExpression: aws.String("#uid = :uid AND #jid = :jid"),
		ExpressionAttributeNames: map[string]string{
			"#uid": constants.DynamoDbKeyUserId,
			"#jid": constants.DynamoDbKeyJournalId,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userId},
			":jid": &types.AttributeValueMemberS{Value: journalId},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query GSI: %w", err)
	}
	if len(result.Items) == 0 {
		return nil, fmt.Errorf("journal entry not found")
	}
	var journal models.Journal
	err = attributevalue.UnmarshalMap(result.Items[0], &journal)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal journal: %w", err)
	}
	return &journal, nil
}

// UpdateJournalEntry updates an existing journal entry using the GSI to find createdAt
func UpdateJournalEntry(ctx context.Context, userId string, journalId string, title, content string) error {
	journal, err := GetJournalByID(ctx, userId, journalId)
	if err != nil {
		return err
	}
	key := map[string]types.AttributeValue{
		constants.DynamoDbKeyUserId: &types.AttributeValueMemberS{Value: userId},
		"CreatedAt":                 &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", journal.CreatedAt)},
	}

	updateExpression := "SET #title = :title, #content = :content, #updatedAt = :updatedAt"
	expressionAttributeNames := map[string]string{
		"#title":     "Title",
		"#content":   "Content",
		"#updatedAt": "UpdatedAt",
	}
	expressionAttributeValues := map[string]types.AttributeValue{
		":title":     &types.AttributeValueMemberS{Value: title},
		":content":   &types.AttributeValueMemberS{Value: content},
		":updatedAt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
	}

	_, err = GetInitializedClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(constants.JournalsTable),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	})
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// DeleteJournalEntry deletes a journal entry from DynamoDB using the GSI to find createdAt
func DeleteJournalEntry(ctx context.Context, userId string, journalId string) error {
	journal, err := GetJournalByID(ctx, userId, journalId)
	if err != nil {
		return err
	}
	key := map[string]types.AttributeValue{
		constants.DynamoDbKeyUserId: &types.AttributeValueMemberS{Value: userId},
		"CreatedAt":                 &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", journal.CreatedAt)},
	}

	_, err = GetInitializedClient().DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(constants.JournalsTable),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// GetUserJournals retrieves up to a limit of journal entries for a user, ordered by most recent (reverse chronological)
func GetUserJournals(ctx context.Context, userId string) ([]models.Journal, error) {
	limit := constants.JournalQueryLimit
	result, err := GetInitializedClient().Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(constants.JournalsTable),
		KeyConditionExpression: aws.String("#uid = :uid"),
		ExpressionAttributeNames: map[string]string{
			"#uid": constants.DynamoDbKeyUserId,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userId},
		},
		ScanIndexForward: aws.Bool(false), // descending order
		Limit:            aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}

	var journals []models.Journal
	for _, item := range result.Items {
		var journal models.Journal
		err := attributevalue.UnmarshalMap(item, &journal)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		journals = append(journals, journal)
	}

	return journals, nil
}
