package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PutAliasRequestAdditionalProperties struct {
	ManagementGroupId    *string            `json:"managementGroupId,omitempty"`
	SubscriptionOwnerId  *string            `json:"subscriptionOwnerId,omitempty"`
	SubscriptionTenantId *string            `json:"subscriptionTenantId,omitempty"`
	Tags                 *map[string]string `json:"tags,omitempty"`
}
