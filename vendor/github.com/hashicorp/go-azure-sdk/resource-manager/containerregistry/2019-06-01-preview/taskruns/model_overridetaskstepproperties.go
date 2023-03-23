package taskruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OverrideTaskStepProperties struct {
	Arguments          *[]Argument `json:"arguments,omitempty"`
	ContextPath        *string     `json:"contextPath,omitempty"`
	File               *string     `json:"file,omitempty"`
	Target             *string     `json:"target,omitempty"`
	UpdateTriggerToken *string     `json:"updateTriggerToken,omitempty"`
	Values             *[]SetValue `json:"values,omitempty"`
}
