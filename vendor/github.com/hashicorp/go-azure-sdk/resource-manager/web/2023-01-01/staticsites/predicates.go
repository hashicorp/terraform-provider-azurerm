package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseConnectionOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p DatabaseConnectionOperationPredicate) Matches(input DatabaseConnection) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type RemotePrivateEndpointConnectionARMResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p RemotePrivateEndpointConnectionARMResourceOperationPredicate) Matches(input RemotePrivateEndpointConnectionARMResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type StaticSiteARMResourceOperationPredicate struct {
	Id       *string
	Kind     *string
	Location *string
	Name     *string
	Type     *string
}

func (p StaticSiteARMResourceOperationPredicate) Matches(input StaticSiteARMResource) bool {

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

type StaticSiteBasicAuthPropertiesARMResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p StaticSiteBasicAuthPropertiesARMResourceOperationPredicate) Matches(input StaticSiteBasicAuthPropertiesARMResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type StaticSiteBuildARMResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p StaticSiteBuildARMResourceOperationPredicate) Matches(input StaticSiteBuildARMResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type StaticSiteCustomDomainOverviewARMResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p StaticSiteCustomDomainOverviewARMResourceOperationPredicate) Matches(input StaticSiteCustomDomainOverviewARMResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type StaticSiteFunctionOverviewARMResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p StaticSiteFunctionOverviewARMResourceOperationPredicate) Matches(input StaticSiteFunctionOverviewARMResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type StaticSiteLinkedBackendARMResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p StaticSiteLinkedBackendARMResourceOperationPredicate) Matches(input StaticSiteLinkedBackendARMResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type StaticSiteUserARMResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p StaticSiteUserARMResourceOperationPredicate) Matches(input StaticSiteUserARMResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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

type StaticSiteUserProvidedFunctionAppARMResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p StaticSiteUserProvidedFunctionAppARMResourceOperationPredicate) Matches(input StaticSiteUserProvidedFunctionAppARMResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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
