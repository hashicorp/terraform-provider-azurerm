package registrationdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationDefinitionProperties struct {
	Authorizations             []Authorization    `json:"authorizations"`
	Description                *string            `json:"description,omitempty"`
	ManagedByTenantId          string             `json:"managedByTenantId"`
	ManagedByTenantName        *string            `json:"managedByTenantName,omitempty"`
	ProvisioningState          *ProvisioningState `json:"provisioningState,omitempty"`
	RegistrationDefinitionName *string            `json:"registrationDefinitionName,omitempty"`
}
