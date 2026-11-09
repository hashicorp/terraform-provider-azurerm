package expressroutelagmembers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagMemberOperationPredicate struct {
	Etag *string
	Id   *string
	Name *string
	Type *string
}

func (p ExpressRouteLagMemberOperationPredicate) Matches(input ExpressRouteLagMember) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

	if p.Id != nil && *p.Id != input.Id {
		return false
	}

	if p.Name != nil && *p.Name != input.Name {
		return false
	}

	if p.Type != nil && *p.Type != input.Type {
		return false
	}

	return true
}
