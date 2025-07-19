package database

import (
	"context"
	"lambda-server/constants"
	"lambda-server/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// CreateMindMuseScoreEntry inserts a new score entry into the mindMuse_score table
func CreateMindMuseScoreEntry(ctx context.Context, score models.MindMuseScore) error {
	client := GetInitializedClient()
	item, err := attributevalue.MarshalMap(score)
	if err != nil {
		return err
	}
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(constants.MindMuseScoreTable),
		Item:      item,
	})
	return err
} 