package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AcceptOwnershipStatusResponse struct {
	AcceptOwnershipState *AcceptOwnership   `json:"acceptOwnershipState,omitempty"`
	BillingOwner         *string            `json:"billingOwner,omitempty"`
	DisplayName          *string            `json:"displayName,omitempty"`
	ProvisioningState    *Provisioning      `json:"provisioningState,omitempty"`
	SubscriptionId       *string            `json:"subscriptionId,omitempty"`
	SubscriptionTenantId *string            `json:"subscriptionTenantId,omitempty"`
	Tags                 *map[string]string `json:"tags,omitempty"`
}
