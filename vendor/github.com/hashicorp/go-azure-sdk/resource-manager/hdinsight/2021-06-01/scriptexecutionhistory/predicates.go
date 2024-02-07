package scriptexecutionhistory

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuntimeScriptActionDetailOperationPredicate struct {
	ApplicationName   *string
	DebugInformation  *string
	EndTime           *string
	Name              *string
	Operation         *string
	Parameters        *string
	ScriptExecutionId *int64
	StartTime         *string
	Status            *string
	Uri               *string
}

func (p RuntimeScriptActionDetailOperationPredicate) Matches(input RuntimeScriptActionDetail) bool {

	if p.ApplicationName != nil && (input.ApplicationName == nil || *p.ApplicationName != *input.ApplicationName) {
		return false
	}

	if p.DebugInformation != nil && (input.DebugInformation == nil || *p.DebugInformation != *input.DebugInformation) {
		return false
	}

	if p.EndTime != nil && (input.EndTime == nil || *p.EndTime != *input.EndTime) {
		return false
	}

	if p.Name != nil && *p.Name != input.Name {
		return false
	}

	if p.Operation != nil && (input.Operation == nil || *p.Operation != *input.Operation) {
		return false
	}

	if p.Parameters != nil && (input.Parameters == nil || *p.Parameters != *input.Parameters) {
		return false
	}

	if p.ScriptExecutionId != nil && (input.ScriptExecutionId == nil || *p.ScriptExecutionId != *input.ScriptExecutionId) {
		return false
	}

	if p.StartTime != nil && (input.StartTime == nil || *p.StartTime != *input.StartTime) {
		return false
	}

	if p.Status != nil && (input.Status == nil || *p.Status != *input.Status) {
		return false
	}

	if p.Uri != nil && *p.Uri != input.Uri {
		return false
	}

	return true
}
