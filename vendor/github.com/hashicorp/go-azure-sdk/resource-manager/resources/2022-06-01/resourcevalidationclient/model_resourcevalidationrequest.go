package resourcevalidationclient

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceValidationRequest struct {
	Location                              *string                             `json:"location,omitempty"`
	PerformPreflightWithoutRbacWriteCheck *bool                               `json:"performPreflightWithoutRbacWriteCheck,omitempty"`
	Provider                              string                              `json:"provider"`
	Resources                             []ResourceValidationRequestResource `json:"resources"`
	Scope                                 string                              `json:"scope"`
	Type                                  string                              `json:"type"`
	ValidationType                        *ResourceValidationType             `json:"validationType,omitempty"`
}
