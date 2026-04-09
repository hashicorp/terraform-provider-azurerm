package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Peers struct {
	AvailabilityZone *string `json:"availabilityZone,omitempty"`
	SubscriptionId   *string `json:"subscriptionId,omitempty"`
}
