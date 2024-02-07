package capabilities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapabilityProperties struct {
	Description      *string `json:"description,omitempty"`
	ParametersSchema *string `json:"parametersSchema,omitempty"`
	Publisher        *string `json:"publisher,omitempty"`
	TargetType       *string `json:"targetType,omitempty"`
	Urn              *string `json:"urn,omitempty"`
}
