package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableServiceSkuOperationPredicate struct {
	ResourceType *string
}

func (p AvailableServiceSkuOperationPredicate) Matches(input AvailableServiceSku) bool {

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}

type DataMigrationServiceOperationPredicate struct {
	Etag     *string
	Id       *string
	Kind     *string
	Location *string
	Name     *string
	Type     *string
}

func (p DataMigrationServiceOperationPredicate) Matches(input DataMigrationService) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type ProjectTaskOperationPredicate struct {
	Etag *string
	Id   *string
	Name *string
	Type *string
}

func (p ProjectTaskOperationPredicate) Matches(input ProjectTask) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

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
