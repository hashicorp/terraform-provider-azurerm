package bastionhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionActiveSessionListResultOperationPredicate struct {
	NextLink *string
}

func (p BastionActiveSessionListResultOperationPredicate) Matches(input BastionActiveSessionListResult) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}

type BastionHostOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p BastionHostOperationPredicate) Matches(input BastionHost) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
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

type BastionSessionDeleteResultOperationPredicate struct {
	NextLink *string
}

func (p BastionSessionDeleteResultOperationPredicate) Matches(input BastionSessionDeleteResult) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}

type BastionShareableLinkListResultOperationPredicate struct {
	NextLink *string
}

func (p BastionShareableLinkListResultOperationPredicate) Matches(input BastionShareableLinkListResult) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}
