package fabriccapacities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricCapacityOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p FabricCapacityOperationPredicate) Matches(input FabricCapacity) bool {

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

type RpSkuDetailsForExistingResourceOperationPredicate struct {
	ResourceType *string
}

func (p RpSkuDetailsForExistingResourceOperationPredicate) Matches(input RpSkuDetailsForExistingResource) bool {

	if p.ResourceType != nil && *p.ResourceType != input.ResourceType {
		return false
	}

	return true
}

type RpSkuDetailsForNewResourceOperationPredicate struct {
	Name         *string
	ResourceType *string
}

func (p RpSkuDetailsForNewResourceOperationPredicate) Matches(input RpSkuDetailsForNewResource) bool {

	if p.Name != nil && *p.Name != input.Name {
		return false
	}

	if p.ResourceType != nil && *p.ResourceType != input.ResourceType {
		return false
	}

	return true
}
