package ascusages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceUsageOperationPredicate struct {
	CurrentValue *int64
	Limit        *int64
	Unit         *string
}

func (p ResourceUsageOperationPredicate) Matches(input ResourceUsage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil || *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	if p.Unit != nil && (input.Unit == nil || *p.Unit != *input.Unit) {
		return false
	}

	return true
}
