package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SetVariableActivityTypeProperties struct {
	SetSystemVariable *bool        `json:"setSystemVariable,omitempty"`
	Value             *interface{} `json:"value,omitempty"`
	VariableName      *string      `json:"variableName,omitempty"`
}
