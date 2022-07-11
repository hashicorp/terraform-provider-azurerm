package policyinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ErrorDefinition struct {
	AdditionalInfo *[]TypedErrorInfo  `json:"additionalInfo,omitempty"`
	Code           *string            `json:"code,omitempty"`
	Details        *[]ErrorDefinition `json:"details,omitempty"`
	Message        *string            `json:"message,omitempty"`
	Target         *string            `json:"target,omitempty"`
}
