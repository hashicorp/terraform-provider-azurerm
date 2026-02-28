package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriberQueueLimit struct {
	Length   *int64                         `json:"length,omitempty"`
	Strategy *SubscriberMessageDropStrategy `json:"strategy,omitempty"`
}
