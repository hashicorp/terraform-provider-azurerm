package remediations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RemediationOperationPredicate) Matches(input Remediation) bool {

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

type RemediationDeploymentsListResultOperationPredicate struct {
	NextLink *string
}

func (p RemediationDeploymentsListResultOperationPredicate) Matches(input RemediationDeploymentsListResult) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}
