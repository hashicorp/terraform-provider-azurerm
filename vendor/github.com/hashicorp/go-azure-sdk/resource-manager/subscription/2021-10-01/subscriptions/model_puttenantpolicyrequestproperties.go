package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PutTenantPolicyRequestProperties struct {
	BlockSubscriptionsIntoTenant    *bool     `json:"blockSubscriptionsIntoTenant,omitempty"`
	BlockSubscriptionsLeavingTenant *bool     `json:"blockSubscriptionsLeavingTenant,omitempty"`
	ExemptedPrincipals              *[]string `json:"exemptedPrincipals,omitempty"`
}
