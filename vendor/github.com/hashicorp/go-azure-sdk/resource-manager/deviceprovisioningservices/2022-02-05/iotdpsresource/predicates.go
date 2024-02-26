package iotdpsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotDpsSkuDefinitionOperationPredicate struct {
}

func (p IotDpsSkuDefinitionOperationPredicate) Matches(input IotDpsSkuDefinition) bool {

	return true
}

type ProvisioningServiceDescriptionOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ProvisioningServiceDescriptionOperationPredicate) Matches(input ProvisioningServiceDescription) bool {

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

type SharedAccessSignatureAuthorizationRuleListResultOperationPredicate struct {
	NextLink *string
}

func (p SharedAccessSignatureAuthorizationRuleListResultOperationPredicate) Matches(input SharedAccessSignatureAuthorizationRuleListResult) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}
