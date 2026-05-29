package linkedresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedResourceOperationPredicate struct {
	Id *string
}

func (p LinkedResourceOperationPredicate) Matches(input LinkedResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	return true
}
