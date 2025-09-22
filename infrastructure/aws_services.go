package infrastructure

import (
	/*
	"context"
	"time"
	*/

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joy-currency-conversion-private/domain"
)

// AWSServices contains all AWS service clients and implementations
type AWSServices struct {
	// AWS SDK clients
	DynamoDB *dynamodb.DynamoDB
	SES      *ses.SES
	SQS      *sqs.SQS

	// Service implementations
	CurrencyService    domain.CurrencyService
	FavoriteService    domain.FavoriteService
	NotificationService domain.NotificationService
}

// NewAWSServices creates a new AWSServices instance
func NewAWSServices(exchangeRateAPIKey string) *AWSServices {
	// Create AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Configure your preferred region
	}))

	// Initialize AWS clients
	dynamoDB := dynamodb.New(sess)
	sesClient := ses.New(sess)
	sqsClient := sqs.New(sess)

	// Initialize service implementations
	currencyService := NewCurrencyService(dynamoDB, exchangeRateAPIKey)
	favoriteService := NewFavoriteService(dynamoDB, exchangeRateAPIKey)
	notificationService := NewNotificationService(sesClient, sqsClient)

	return &AWSServices{
		DynamoDB:            dynamoDB,
		SES:                 sesClient,
		SQS:                 sqsClient,
		CurrencyService:     currencyService,
		FavoriteService:     favoriteService,
		NotificationService: notificationService,
	}
}
