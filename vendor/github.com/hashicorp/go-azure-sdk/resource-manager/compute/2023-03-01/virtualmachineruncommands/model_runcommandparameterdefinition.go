package virtualmachineruncommands

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunCommandParameterDefinition struct {
	DefaultValue *string `json:"defaultValue,omitempty"`
	Name         string  `json:"name"`
	Required     *bool   `json:"required,omitempty"`
	Type         string  `json:"type"`
}
