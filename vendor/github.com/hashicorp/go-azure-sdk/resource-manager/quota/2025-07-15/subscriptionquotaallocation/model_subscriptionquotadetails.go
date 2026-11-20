package subscriptionquotaallocation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionQuotaDetails struct {
	Limit          *int64                        `json:"limit,omitempty"`
	Name           *SubscriptionQuotaDetailsName `json:"name,omitempty"`
	ResourceName   *string                       `json:"resourceName,omitempty"`
	ShareableQuota *int64                        `json:"shareableQuota,omitempty"`
}
