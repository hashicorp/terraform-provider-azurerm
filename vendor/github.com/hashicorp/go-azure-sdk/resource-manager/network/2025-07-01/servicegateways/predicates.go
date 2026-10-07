package servicegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceGatewayOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ServiceGatewayOperationPredicate) Matches(input ServiceGateway) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

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

type ServiceGatewayAddressLocationResponseOperationPredicate struct {
	AddressLocation *string
}

func (p ServiceGatewayAddressLocationResponseOperationPredicate) Matches(input ServiceGatewayAddressLocationResponse) bool {

	if p.AddressLocation != nil && (input.AddressLocation == nil || *p.AddressLocation != *input.AddressLocation) {
		return false
	}

	return true
}

type ServiceGatewayServiceOperationPredicate struct {
	Name *string
}

func (p ServiceGatewayServiceOperationPredicate) Matches(input ServiceGatewayService) bool {

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	return true
}
