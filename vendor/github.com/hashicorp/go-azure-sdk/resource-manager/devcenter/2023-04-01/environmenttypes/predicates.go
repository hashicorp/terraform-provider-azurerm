package environmenttypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllowedEnvironmentTypeOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p AllowedEnvironmentTypeOperationPredicate) Matches(input AllowedEnvironmentType) bool {

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

type EnvironmentTypeOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p EnvironmentTypeOperationPredicate) Matches(input EnvironmentType) bool {

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

type ProjectEnvironmentTypeOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ProjectEnvironmentTypeOperationPredicate) Matches(input ProjectEnvironmentType) bool {

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
