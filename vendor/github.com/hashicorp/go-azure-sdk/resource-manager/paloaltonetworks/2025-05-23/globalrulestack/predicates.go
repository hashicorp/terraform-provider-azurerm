package globalrulestack

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvSecurityObjectModelOperationPredicate struct {
	Type *string
}

func (p AdvSecurityObjectModelOperationPredicate) Matches(input AdvSecurityObjectModel) bool {

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type CountryOperationPredicate struct {
	Code        *string
	Description *string
}

func (p CountryOperationPredicate) Matches(input Country) bool {

	if p.Code != nil && *p.Code != input.Code {
		return false
	}

	if p.Description != nil && (input.Description == nil || *p.Description != *input.Description) {
		return false
	}

	return true
}

type GlobalRulestackResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p GlobalRulestackResourceOperationPredicate) Matches(input GlobalRulestackResource) bool {

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

type PredefinedURLCategoryOperationPredicate struct {
	Action *string
	Name   *string
}

func (p PredefinedURLCategoryOperationPredicate) Matches(input PredefinedURLCategory) bool {

	if p.Action != nil && *p.Action != input.Action {
		return false
	}

	if p.Name != nil && *p.Name != input.Name {
		return false
	}

	return true
}

type SecurityServicesTypeListOperationPredicate struct {
	Type *string
}

func (p SecurityServicesTypeListOperationPredicate) Matches(input SecurityServicesTypeList) bool {

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}
