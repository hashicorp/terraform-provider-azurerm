package eventsubscriptions

type StorageQueueEventSubscriptionDestinationProperties struct {
	QueueMessageTimeToLiveInSeconds *int64  `json:"queueMessageTimeToLiveInSeconds,omitempty"`
	QueueName                       *string `json:"queueName,omitempty"`
	ResourceId                      *string `json:"resourceId,omitempty"`
}
