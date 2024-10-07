package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CachedImagesOperationPredicate struct {
	Image  *string
	OsType *string
}

func (p CachedImagesOperationPredicate) Matches(input CachedImages) bool {

	if p.Image != nil && *p.Image != input.Image {
		return false
	}

	if p.OsType != nil && *p.OsType != input.OsType {
		return false
	}

	return true
}

type CapabilitiesOperationPredicate struct {
	Gpu           *string
	IPAddressType *string
	Location      *string
	OsType        *string
	ResourceType  *string
}

func (p CapabilitiesOperationPredicate) Matches(input Capabilities) bool {

	if p.Gpu != nil && (input.Gpu == nil || *p.Gpu != *input.Gpu) {
		return false
	}

	if p.IPAddressType != nil && (input.IPAddressType == nil || *p.IPAddressType != *input.IPAddressType) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.OsType != nil && (input.OsType == nil || *p.OsType != *input.OsType) {
		return false
	}

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}

type ContainerGroupOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ContainerGroupOperationPredicate) Matches(input ContainerGroup) bool {

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
