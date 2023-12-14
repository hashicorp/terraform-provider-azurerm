package messages

import (
	"context"
)

type StorageQueueMessage interface {
	Delete(ctx context.Context, queueName string, messageID string, input DeleteInput) (DeleteResponse, error)
	Peek(ctx context.Context, queueName string, input PeekInput) (QueueMessagesListResponse, error)
	Put(ctx context.Context, queueName string, input PutInput) (QueueMessagesListResponse, error)
	Get(ctx context.Context, queueName string, input GetInput) (QueueMessagesListResponse, error)
	Update(ctx context.Context, queueName string, messageID string, input UpdateInput) (UpdateResponse, error)
}
