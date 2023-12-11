package eventsubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetryPolicy struct {
	EventTimeToLiveInMinutes *int64 `json:"eventTimeToLiveInMinutes,omitempty"`
	MaxDeliveryAttempts      *int64 `json:"maxDeliveryAttempts,omitempty"`
}
