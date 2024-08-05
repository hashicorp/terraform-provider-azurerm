package globalrulestack

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RulestackProperties struct {
	AssociatedSubscriptions *[]string          `json:"associatedSubscriptions,omitempty"`
	DefaultMode             *DefaultMode       `json:"defaultMode,omitempty"`
	Description             *string            `json:"description,omitempty"`
	MinAppIdVersion         *string            `json:"minAppIdVersion,omitempty"`
	PanEtag                 *string            `json:"panEtag,omitempty"`
	PanLocation             *string            `json:"panLocation,omitempty"`
	ProvisioningState       *ProvisioningState `json:"provisioningState,omitempty"`
	Scope                   *ScopeType         `json:"scope,omitempty"`
	SecurityServices        *SecurityServices  `json:"securityServices,omitempty"`
}
