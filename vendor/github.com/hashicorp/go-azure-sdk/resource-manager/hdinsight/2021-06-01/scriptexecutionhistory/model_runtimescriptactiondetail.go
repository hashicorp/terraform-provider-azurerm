package scriptexecutionhistory

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuntimeScriptActionDetail struct {
	ApplicationName   *string                         `json:"applicationName,omitempty"`
	DebugInformation  *string                         `json:"debugInformation,omitempty"`
	EndTime           *string                         `json:"endTime,omitempty"`
	ExecutionSummary  *[]ScriptActionExecutionSummary `json:"executionSummary,omitempty"`
	Name              string                          `json:"name"`
	Operation         *string                         `json:"operation,omitempty"`
	Parameters        *string                         `json:"parameters,omitempty"`
	Roles             []string                        `json:"roles"`
	ScriptExecutionId *int64                          `json:"scriptExecutionId,omitempty"`
	StartTime         *string                         `json:"startTime,omitempty"`
	Status            *string                         `json:"status,omitempty"`
	Uri               string                          `json:"uri"`
}
