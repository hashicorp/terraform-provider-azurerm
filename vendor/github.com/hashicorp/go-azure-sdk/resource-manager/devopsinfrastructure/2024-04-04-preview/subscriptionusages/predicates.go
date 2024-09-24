package subscriptionusages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaOperationPredicate struct {
	CurrentValue *int64
	Id           *string
	Limit        *int64
	Unit         *string
}

func (p QuotaOperationPredicate) Matches(input Quota) bool {

	if p.CurrentValue != nil && *p.CurrentValue != input.CurrentValue {
		return false
	}

	if p.Id != nil && *p.Id != input.Id {
		return false
	}

	if p.Limit != nil && *p.Limit != input.Limit {
		return false
	}

	if p.Unit != nil && *p.Unit != input.Unit {
		return false
	}

	return true
}
