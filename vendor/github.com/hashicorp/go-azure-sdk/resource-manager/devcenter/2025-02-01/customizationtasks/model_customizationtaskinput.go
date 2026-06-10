package customizationtasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomizationTaskInput struct {
	Description *string                     `json:"description,omitempty"`
	Required    *bool                       `json:"required,omitempty"`
	Type        *CustomizationTaskInputType `json:"type,omitempty"`
}
