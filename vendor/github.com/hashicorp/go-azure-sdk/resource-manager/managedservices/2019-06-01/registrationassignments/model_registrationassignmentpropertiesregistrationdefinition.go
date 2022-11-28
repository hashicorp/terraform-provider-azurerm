package registrationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationAssignmentPropertiesRegistrationDefinition struct {
	Id         *string                                                           `json:"id,omitempty"`
	Name       *string                                                           `json:"name,omitempty"`
	Plan       *Plan                                                             `json:"plan"`
	Properties *RegistrationAssignmentPropertiesRegistrationDefinitionProperties `json:"properties"`
	Type       *string                                                           `json:"type,omitempty"`
}
