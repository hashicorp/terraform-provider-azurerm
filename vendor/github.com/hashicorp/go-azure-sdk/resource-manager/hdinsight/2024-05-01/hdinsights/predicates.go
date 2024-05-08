package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ClusterOperationPredicate) Matches(input Cluster) bool {

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

type ClusterAvailableUpgradeOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClusterAvailableUpgradeOperationPredicate) Matches(input ClusterAvailableUpgrade) bool {

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

type ClusterInstanceViewResultOperationPredicate struct {
	Name *string
}

func (p ClusterInstanceViewResultOperationPredicate) Matches(input ClusterInstanceViewResult) bool {

	if p.Name != nil && *p.Name != input.Name {
		return false
	}

	return true
}

type ClusterJobOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClusterJobOperationPredicate) Matches(input ClusterJob) bool {

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

type ClusterLibraryOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClusterLibraryOperationPredicate) Matches(input ClusterLibrary) bool {

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

type ClusterPoolOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ClusterPoolOperationPredicate) Matches(input ClusterPool) bool {

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

type ClusterPoolAvailableUpgradeOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClusterPoolAvailableUpgradeOperationPredicate) Matches(input ClusterPoolAvailableUpgrade) bool {

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

type ClusterPoolUpgradeHistoryOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClusterPoolUpgradeHistoryOperationPredicate) Matches(input ClusterPoolUpgradeHistory) bool {

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

type ClusterPoolVersionOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClusterPoolVersionOperationPredicate) Matches(input ClusterPoolVersion) bool {

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

type ClusterUpgradeHistoryOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClusterUpgradeHistoryOperationPredicate) Matches(input ClusterUpgradeHistory) bool {

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

type ClusterVersionOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClusterVersionOperationPredicate) Matches(input ClusterVersion) bool {

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

type ServiceConfigResultOperationPredicate struct {
}

func (p ServiceConfigResultOperationPredicate) Matches(input ServiceConfigResult) bool {

	return true
}
