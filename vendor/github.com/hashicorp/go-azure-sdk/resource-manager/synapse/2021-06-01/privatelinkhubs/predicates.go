package privatelinkhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionForPrivateLinkHubOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p PrivateEndpointConnectionForPrivateLinkHubOperationPredicate) Matches(input PrivateEndpointConnectionForPrivateLinkHub) bool {

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

type PrivateLinkHubOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p PrivateLinkHubOperationPredicate) Matches(input PrivateLinkHub) bool {

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
