package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PutAliasRequestProperties struct {
	AdditionalProperties *PutAliasRequestAdditionalProperties `json:"additionalProperties,omitempty"`
	BillingScope         *string                              `json:"billingScope,omitempty"`
	DisplayName          *string                              `json:"displayName,omitempty"`
	ResellerId           *string                              `json:"resellerId,omitempty"`
	SubscriptionId       *string                              `json:"subscriptionId,omitempty"`
	Workload             *Workload                            `json:"workload,omitempty"`
}
