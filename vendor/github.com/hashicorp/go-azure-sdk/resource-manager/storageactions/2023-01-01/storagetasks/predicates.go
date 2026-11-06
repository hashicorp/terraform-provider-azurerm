package storagetasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p StorageTaskOperationPredicate) Matches(input StorageTask) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type StorageTaskAssignmentOperationPredicate struct {
	Id *string
}

func (p StorageTaskAssignmentOperationPredicate) Matches(input StorageTaskAssignment) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	return true
}

type StorageTaskReportInstanceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p StorageTaskReportInstanceOperationPredicate) Matches(input StorageTaskReportInstance) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}
