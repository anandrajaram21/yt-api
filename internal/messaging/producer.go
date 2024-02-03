package messaging

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SendMessage(queueURL, messageBody string) error {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return err
	}

	sqsClient := sqs.NewFromConfig(cfg)

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	log.Printf("Message sent to SQS queue: %s", messageBody)
	return nil
}
