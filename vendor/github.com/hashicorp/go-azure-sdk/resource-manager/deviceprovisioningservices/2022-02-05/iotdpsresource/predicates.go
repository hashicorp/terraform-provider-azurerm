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

type SharedAccessSignatureAuthorizationRuleAccessRightsDescriptionOperationPredicate struct {
	KeyName      *string
	PrimaryKey   *string
	SecondaryKey *string
}

func (p SharedAccessSignatureAuthorizationRuleAccessRightsDescriptionOperationPredicate) Matches(input SharedAccessSignatureAuthorizationRuleAccessRightsDescription) bool {

	if p.KeyName != nil && *p.KeyName != input.KeyName {
		return false
	}

	if p.PrimaryKey != nil && (input.PrimaryKey == nil || *p.PrimaryKey != *input.PrimaryKey) {
		return false
	}

	if p.SecondaryKey != nil && (input.SecondaryKey == nil || *p.SecondaryKey != *input.SecondaryKey) {
		return false
	}

	return true
}
