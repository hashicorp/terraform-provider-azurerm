package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableDelegationOperationPredicate struct {
	Id          *string
	Name        *string
	ServiceName *string
	Type        *string
}

func (p AvailableDelegationOperationPredicate) Matches(input AvailableDelegation) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.ServiceName != nil && (input.ServiceName == nil || *p.ServiceName != *input.ServiceName) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type AvailablePrivateEndpointTypeOperationPredicate struct {
	DisplayName  *string
	Id           *string
	Name         *string
	ResourceName *string
	Type         *string
}

func (p AvailablePrivateEndpointTypeOperationPredicate) Matches(input AvailablePrivateEndpointType) bool {

	if p.DisplayName != nil && (input.DisplayName == nil || *p.DisplayName != *input.DisplayName) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.ResourceName != nil && (input.ResourceName == nil || *p.ResourceName != *input.ResourceName) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type AvailableServiceAliasOperationPredicate struct {
	Id           *string
	Name         *string
	ResourceName *string
	Type         *string
}

func (p AvailableServiceAliasOperationPredicate) Matches(input AvailableServiceAlias) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.ResourceName != nil && (input.ResourceName == nil || *p.ResourceName != *input.ResourceName) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type EndpointServiceResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p EndpointServiceResultOperationPredicate) Matches(input EndpointServiceResult) bool {

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

type PublicIPAddressOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p PublicIPAddressOperationPredicate) Matches(input PublicIPAddress) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

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

type PublicIPDdosProtectionStatusResultOperationPredicate struct {
	DdosProtectionPlanId *string
	PublicIPAddress      *string
	PublicIPAddressId    *string
}

func (p PublicIPDdosProtectionStatusResultOperationPredicate) Matches(input PublicIPDdosProtectionStatusResult) bool {

	if p.DdosProtectionPlanId != nil && (input.DdosProtectionPlanId == nil || *p.DdosProtectionPlanId != *input.DdosProtectionPlanId) {
		return false
	}

	if p.PublicIPAddress != nil && (input.PublicIPAddress == nil || *p.PublicIPAddress != *input.PublicIPAddress) {
		return false
	}

	if p.PublicIPAddressId != nil && (input.PublicIPAddressId == nil || *p.PublicIPAddressId != *input.PublicIPAddressId) {
		return false
	}

	return true
}

type ResourceNavigationLinkOperationPredicate struct {
	Etag *string
	Id   *string
	Name *string
	Type *string
}

func (p ResourceNavigationLinkOperationPredicate) Matches(input ResourceNavigationLink) bool {

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

type ServiceAssociationLinkOperationPredicate struct {
	Etag *string
	Id   *string
	Name *string
	Type *string
}

func (p ServiceAssociationLinkOperationPredicate) Matches(input ServiceAssociationLink) bool {

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

type ServiceTagInformationOperationPredicate struct {
	Id                     *string
	Name                   *string
	ServiceTagChangeNumber *string
}

func (p ServiceTagInformationOperationPredicate) Matches(input ServiceTagInformation) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.ServiceTagChangeNumber != nil && (input.ServiceTagChangeNumber == nil || *p.ServiceTagChangeNumber != *input.ServiceTagChangeNumber) {
		return false
	}

	return true
}

type VirtualNetworkOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p VirtualNetworkOperationPredicate) Matches(input VirtualNetwork) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

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

type VirtualNetworkUsageOperationPredicate struct {
	CurrentValue *float64
	Id           *string
	Limit        *float64
	Unit         *string
}

func (p VirtualNetworkUsageOperationPredicate) Matches(input VirtualNetworkUsage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil || *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	if p.Unit != nil && (input.Unit == nil || *p.Unit != *input.Unit) {
		return false
	}

	return true
}
