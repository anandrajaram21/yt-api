package messaging

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func ReceiveMessage(queueURL string) (*types.Message, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	sqsClient := sqs.NewFromConfig(cfg)

	msgOutput, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   30,
	})
	if err != nil {
		log.Printf("Failed to receive messages: %v", err)
		return nil, err
	}

	if len(msgOutput.Messages) == 0 {
		log.Println("No messages received")
		return nil, nil
	}

	message := msgOutput.Messages[0]
	log.Printf("Received message: %s", aws.ToString(message.Body))

	// Acknowledge the message by deleting it from the queue
	_, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
		return &message, err
	}

	log.Printf("Message deleted successfully")

	return &message, nil
}
