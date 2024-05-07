package bastionshareablelink

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionShareableLinkOperationPredicate struct {
	Bsl       *string
	CreatedAt *string
	Message   *string
}

func (p BastionShareableLinkOperationPredicate) Matches(input BastionShareableLink) bool {

	if p.Bsl != nil && (input.Bsl == nil || *p.Bsl != *input.Bsl) {
		return false
	}

	if p.CreatedAt != nil && (input.CreatedAt == nil || *p.CreatedAt != *input.CreatedAt) {
		return false
	}

	if p.Message != nil && (input.Message == nil || *p.Message != *input.Message) {
		return false
	}

	return true
}
