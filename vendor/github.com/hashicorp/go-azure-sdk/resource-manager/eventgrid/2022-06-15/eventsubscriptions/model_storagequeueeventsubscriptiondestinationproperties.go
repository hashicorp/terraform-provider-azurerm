package eventsubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageQueueEventSubscriptionDestinationProperties struct {
	QueueMessageTimeToLiveInSeconds *int64  `json:"queueMessageTimeToLiveInSeconds,omitempty"`
	QueueName                       *string `json:"queueName,omitempty"`
	ResourceId                      *string `json:"resourceId,omitempty"`
}
