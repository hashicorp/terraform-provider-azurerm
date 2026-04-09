package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceEnvironmentResourceOperationPredicate struct {
	Id       *string
	Kind     *string
	Location *string
	Name     *string
	Type     *string
}

func (p AppServiceEnvironmentResourceOperationPredicate) Matches(input AppServiceEnvironmentResource) bool {

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

type AppServicePlanOperationPredicate struct {
	Id       *string
	Kind     *string
	Location *string
	Name     *string
	Type     *string
}

func (p AppServicePlanOperationPredicate) Matches(input AppServicePlan) bool {

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

type CsmUsageQuotaOperationPredicate struct {
	CurrentValue  *int64
	Limit         *int64
	NextResetTime *string
	Unit          *string
}

func (p CsmUsageQuotaOperationPredicate) Matches(input CsmUsageQuota) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil || *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	if p.NextResetTime != nil && (input.NextResetTime == nil || *p.NextResetTime != *input.NextResetTime) {
		return false
	}

	if p.Unit != nil && (input.Unit == nil || *p.Unit != *input.Unit) {
		return false
	}

	return true
}

type InboundEnvironmentEndpointOperationPredicate struct {
	Description *string
}

func (p InboundEnvironmentEndpointOperationPredicate) Matches(input InboundEnvironmentEndpoint) bool {

	if p.Description != nil && (input.Description == nil || *p.Description != *input.Description) {
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

type ResourceMetricDefinitionOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p ResourceMetricDefinitionOperationPredicate) Matches(input ResourceMetricDefinition) bool {

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

type SiteOperationPredicate struct {
	Id       *string
	Kind     *string
	Location *string
	Name     *string
	Type     *string
}

func (p SiteOperationPredicate) Matches(input Site) bool {

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

type SkuInfoOperationPredicate struct {
	ResourceType *string
}

func (p SkuInfoOperationPredicate) Matches(input SkuInfo) bool {

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	return true
}

type StampCapacityOperationPredicate struct {
	AvailableCapacity              *int64
	ExcludeFromCapacityAllocation  *bool
	IsApplicableForAllComputeModes *bool
	IsLinux                        *bool
	Name                           *string
	SiteMode                       *string
	TotalCapacity                  *int64
	Unit                           *string
	WorkerSizeId                   *int64
}

func (p StampCapacityOperationPredicate) Matches(input StampCapacity) bool {

	if p.AvailableCapacity != nil && (input.AvailableCapacity == nil || *p.AvailableCapacity != *input.AvailableCapacity) {
		return false
	}

	if p.ExcludeFromCapacityAllocation != nil && (input.ExcludeFromCapacityAllocation == nil || *p.ExcludeFromCapacityAllocation != *input.ExcludeFromCapacityAllocation) {
		return false
	}

	if p.IsApplicableForAllComputeModes != nil && (input.IsApplicableForAllComputeModes == nil || *p.IsApplicableForAllComputeModes != *input.IsApplicableForAllComputeModes) {
		return false
	}

	if p.IsLinux != nil && (input.IsLinux == nil || *p.IsLinux != *input.IsLinux) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.SiteMode != nil && (input.SiteMode == nil || *p.SiteMode != *input.SiteMode) {
		return false
	}

	if p.TotalCapacity != nil && (input.TotalCapacity == nil || *p.TotalCapacity != *input.TotalCapacity) {
		return false
	}

	if p.Unit != nil && (input.Unit == nil || *p.Unit != *input.Unit) {
		return false
	}

	if p.WorkerSizeId != nil && (input.WorkerSizeId == nil || *p.WorkerSizeId != *input.WorkerSizeId) {
		return false
	}

	return true
}

type UsageOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p UsageOperationPredicate) Matches(input Usage) bool {

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

type WorkerPoolResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p WorkerPoolResourceOperationPredicate) Matches(input WorkerPoolResource) bool {

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
