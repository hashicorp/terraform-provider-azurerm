package registrationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationAssignmentProperties struct {
	ProvisioningState        *ProvisioningState                                      `json:"provisioningState,omitempty"`
	RegistrationDefinition   *RegistrationAssignmentPropertiesRegistrationDefinition `json:"registrationDefinition,omitempty"`
	RegistrationDefinitionId string                                                  `json:"registrationDefinitionId"`
}
