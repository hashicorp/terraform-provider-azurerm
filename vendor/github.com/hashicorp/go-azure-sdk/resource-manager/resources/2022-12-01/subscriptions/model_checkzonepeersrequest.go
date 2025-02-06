package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckZonePeersRequest struct {
	Location        *string   `json:"location,omitempty"`
	SubscriptionIds *[]string `json:"subscriptionIds,omitempty"`
}
