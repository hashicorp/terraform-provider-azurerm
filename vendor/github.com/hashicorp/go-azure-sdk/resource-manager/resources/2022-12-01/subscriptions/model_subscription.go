package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Subscription struct {
	AuthorizationSource  *string               `json:"authorizationSource,omitempty"`
	DisplayName          *string               `json:"displayName,omitempty"`
	Id                   *string               `json:"id,omitempty"`
	ManagedByTenants     *[]ManagedByTenant    `json:"managedByTenants,omitempty"`
	State                *SubscriptionState    `json:"state,omitempty"`
	SubscriptionId       *string               `json:"subscriptionId,omitempty"`
	SubscriptionPolicies *SubscriptionPolicies `json:"subscriptionPolicies,omitempty"`
	Tags                 *map[string]string    `json:"tags,omitempty"`
	TenantId             *string               `json:"tenantId,omitempty"`
}
