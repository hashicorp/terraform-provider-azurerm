package groupquotalimits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllocatedToSubscription struct {
	QuotaAllocated *int64  `json:"quotaAllocated,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
}
