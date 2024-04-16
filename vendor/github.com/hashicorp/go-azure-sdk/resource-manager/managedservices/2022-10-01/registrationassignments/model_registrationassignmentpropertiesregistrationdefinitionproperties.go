package registrationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationAssignmentPropertiesRegistrationDefinitionProperties struct {
	Authorizations             *[]Authorization         `json:"authorizations,omitempty"`
	Description                *string                  `json:"description,omitempty"`
	EligibleAuthorizations     *[]EligibleAuthorization `json:"eligibleAuthorizations,omitempty"`
	ManagedByTenantId          *string                  `json:"managedByTenantId,omitempty"`
	ManagedByTenantName        *string                  `json:"managedByTenantName,omitempty"`
	ManageeTenantId            *string                  `json:"manageeTenantId,omitempty"`
	ManageeTenantName          *string                  `json:"manageeTenantName,omitempty"`
	ProvisioningState          *ProvisioningState       `json:"provisioningState,omitempty"`
	RegistrationDefinitionName *string                  `json:"registrationDefinitionName,omitempty"`
}
