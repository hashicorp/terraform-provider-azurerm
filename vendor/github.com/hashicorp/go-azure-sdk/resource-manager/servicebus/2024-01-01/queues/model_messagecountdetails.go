package queues

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MessageCountDetails struct {
	ActiveMessageCount             *int64 `json:"activeMessageCount,omitempty"`
	DeadLetterMessageCount         *int64 `json:"deadLetterMessageCount,omitempty"`
	ScheduledMessageCount          *int64 `json:"scheduledMessageCount,omitempty"`
	TransferDeadLetterMessageCount *int64 `json:"transferDeadLetterMessageCount,omitempty"`
	TransferMessageCount           *int64 `json:"transferMessageCount,omitempty"`
}
