package useridentity

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserIdentityContractOperationPredicate struct {
	Id       *string
	Provider *string
}

func (p UserIdentityContractOperationPredicate) Matches(input UserIdentityContract) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Provider != nil && (input.Provider == nil || *p.Provider != *input.Provider) {
		return false
	}

	return true
}
