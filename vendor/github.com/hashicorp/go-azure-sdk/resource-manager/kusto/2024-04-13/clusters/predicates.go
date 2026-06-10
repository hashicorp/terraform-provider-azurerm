package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CalloutPolicyOperationPredicate struct {
	CalloutId       *string
	CalloutUriRegex *string
}

func (p CalloutPolicyOperationPredicate) Matches(input CalloutPolicy) bool {

	if p.CalloutId != nil && (input.CalloutId == nil || *p.CalloutId != *input.CalloutId) {
		return false
	}

	if p.CalloutUriRegex != nil && (input.CalloutUriRegex == nil || *p.CalloutUriRegex != *input.CalloutUriRegex) {
		return false
	}

	return true
}

type FollowerDatabaseDefinitionGetOperationPredicate struct {
}

func (p FollowerDatabaseDefinitionGetOperationPredicate) Matches(input FollowerDatabaseDefinitionGet) bool {

	return true
}
