package policyassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NonComplianceMessage struct {
	Message                     string  `json:"message"`
	PolicyDefinitionReferenceId *string `json:"policyDefinitionReferenceId,omitempty"`
}
