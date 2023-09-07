package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlComputeNodesInformationOperationPredicate struct {
	NextLink *string
}

func (p AmlComputeNodesInformationOperationPredicate) Matches(input AmlComputeNodesInformation) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}

type ComputeResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ComputeResourceOperationPredicate) Matches(input ComputeResource) bool {

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
