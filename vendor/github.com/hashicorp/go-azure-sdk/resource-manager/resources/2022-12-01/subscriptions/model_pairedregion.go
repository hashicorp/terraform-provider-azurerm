package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PairedRegion struct {
	Id             *string `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
}
