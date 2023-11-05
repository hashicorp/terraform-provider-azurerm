package apioperation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterExampleContract struct {
	Description   *string      `json:"description,omitempty"`
	ExternalValue *string      `json:"externalValue,omitempty"`
	Summary       *string      `json:"summary,omitempty"`
	Value         *interface{} `json:"value,omitempty"`
}
