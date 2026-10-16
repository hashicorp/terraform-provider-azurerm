package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoHealActions struct {
	ActionType              *AutoHealActionType   `json:"actionType,omitempty"`
	CustomAction            *AutoHealCustomAction `json:"customAction,omitempty"`
	MinProcessExecutionTime *string               `json:"minProcessExecutionTime,omitempty"`
}
