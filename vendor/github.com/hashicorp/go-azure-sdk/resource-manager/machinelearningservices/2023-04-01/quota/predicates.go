package quota

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceQuotaOperationPredicate struct {
	AmlWorkspaceLocation *string
	Id                   *string
	Limit                *int64
	Type                 *string
}

func (p ResourceQuotaOperationPredicate) Matches(input ResourceQuota) bool {

	if p.AmlWorkspaceLocation != nil && (input.AmlWorkspaceLocation == nil || *p.AmlWorkspaceLocation != *input.AmlWorkspaceLocation) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}
