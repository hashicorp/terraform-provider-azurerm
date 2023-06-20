package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionAliasResponseProperties struct {
	AcceptOwnershipState *AcceptOwnership   `json:"acceptOwnershipState,omitempty"`
	AcceptOwnershipUrl   *string            `json:"acceptOwnershipUrl,omitempty"`
	BillingScope         *string            `json:"billingScope,omitempty"`
	CreatedTime          *string            `json:"createdTime,omitempty"`
	DisplayName          *string            `json:"displayName,omitempty"`
	ManagementGroupId    *string            `json:"managementGroupId,omitempty"`
	ProvisioningState    *ProvisioningState `json:"provisioningState,omitempty"`
	ResellerId           *string            `json:"resellerId,omitempty"`
	SubscriptionId       *string            `json:"subscriptionId,omitempty"`
	SubscriptionOwnerId  *string            `json:"subscriptionOwnerId,omitempty"`
	Tags                 *map[string]string `json:"tags,omitempty"`
	Workload             *Workload          `json:"workload,omitempty"`
}
