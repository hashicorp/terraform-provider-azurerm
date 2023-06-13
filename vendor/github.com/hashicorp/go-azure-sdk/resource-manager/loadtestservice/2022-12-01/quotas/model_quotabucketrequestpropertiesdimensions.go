package quotas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaBucketRequestPropertiesDimensions struct {
	Location       *string `json:"location,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
}
