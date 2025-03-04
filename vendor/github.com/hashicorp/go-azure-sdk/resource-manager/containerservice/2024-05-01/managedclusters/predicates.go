package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ManagedClusterOperationPredicate) Matches(input ManagedCluster) bool {

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

type MeshRevisionProfileOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p MeshRevisionProfileOperationPredicate) Matches(input MeshRevisionProfile) bool {

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

type MeshUpgradeProfileOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p MeshUpgradeProfileOperationPredicate) Matches(input MeshUpgradeProfile) bool {

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

type OutboundEnvironmentEndpointOperationPredicate struct {
	Category *string
}

func (p OutboundEnvironmentEndpointOperationPredicate) Matches(input OutboundEnvironmentEndpoint) bool {

	if p.Category != nil && (input.Category == nil || *p.Category != *input.Category) {
		return false
	}

	return true
}
