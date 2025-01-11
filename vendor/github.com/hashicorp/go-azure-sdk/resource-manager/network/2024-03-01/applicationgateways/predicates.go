package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ApplicationGatewayOperationPredicate) Matches(input ApplicationGateway) bool {

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

type ApplicationGatewaySslPredefinedPolicyOperationPredicate struct {
	Id   *string
	Name *string
}

func (p ApplicationGatewaySslPredefinedPolicyOperationPredicate) Matches(input ApplicationGatewaySslPredefinedPolicy) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	return true
}
