package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstructionSteps struct {
	Description  *string                          `json:"description,omitempty"`
	Instructions *[]ConnectorInstructionModelBase `json:"instructions,omitempty"`
	Title        *string                          `json:"title,omitempty"`
}
