package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientConfig struct {
	MaxKeepAliveSeconds     *int64                `json:"maxKeepAliveSeconds,omitempty"`
	MaxMessageExpirySeconds *int64                `json:"maxMessageExpirySeconds,omitempty"`
	MaxPacketSizeBytes      *int64                `json:"maxPacketSizeBytes,omitempty"`
	MaxReceiveMaximum       *int64                `json:"maxReceiveMaximum,omitempty"`
	MaxSessionExpirySeconds *int64                `json:"maxSessionExpirySeconds,omitempty"`
	SubscriberQueueLimit    *SubscriberQueueLimit `json:"subscriberQueueLimit,omitempty"`
}
