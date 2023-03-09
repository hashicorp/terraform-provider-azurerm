package registrationdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JustInTimeAccessPolicy struct {
	ManagedByTenantApprovers  *[]EligibleApprover     `json:"managedByTenantApprovers,omitempty"`
	MaximumActivationDuration *string                 `json:"maximumActivationDuration,omitempty"`
	MultiFactorAuthProvider   MultiFactorAuthProvider `json:"multiFactorAuthProvider"`
}
