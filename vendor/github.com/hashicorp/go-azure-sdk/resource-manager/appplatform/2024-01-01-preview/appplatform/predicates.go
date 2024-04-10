package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiPortalCustomDomainResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ApiPortalCustomDomainResourceOperationPredicate) Matches(input ApiPortalCustomDomainResource) bool {

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

type ApiPortalResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ApiPortalResourceOperationPredicate) Matches(input ApiPortalResource) bool {

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

type ApmResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ApmResourceOperationPredicate) Matches(input ApmResource) bool {

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

type AppResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p AppResourceOperationPredicate) Matches(input AppResource) bool {

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

type ApplicationAcceleratorResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ApplicationAcceleratorResourceOperationPredicate) Matches(input ApplicationAcceleratorResource) bool {

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

type ApplicationLiveViewResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ApplicationLiveViewResourceOperationPredicate) Matches(input ApplicationLiveViewResource) bool {

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

type BindingResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p BindingResourceOperationPredicate) Matches(input BindingResource) bool {

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

type BuildOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p BuildOperationPredicate) Matches(input Build) bool {

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

type BuildResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p BuildResultOperationPredicate) Matches(input BuildResult) bool {

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

type BuildServiceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p BuildServiceOperationPredicate) Matches(input BuildService) bool {

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

type BuildServiceAgentPoolResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p BuildServiceAgentPoolResourceOperationPredicate) Matches(input BuildServiceAgentPoolResource) bool {

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

type BuilderResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p BuilderResourceOperationPredicate) Matches(input BuilderResource) bool {

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

type BuildpackBindingResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p BuildpackBindingResourceOperationPredicate) Matches(input BuildpackBindingResource) bool {

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

type CertificateResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p CertificateResourceOperationPredicate) Matches(input CertificateResource) bool {

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

type ConfigurationServiceResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ConfigurationServiceResourceOperationPredicate) Matches(input ConfigurationServiceResource) bool {

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

type ContainerRegistryResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ContainerRegistryResourceOperationPredicate) Matches(input ContainerRegistryResource) bool {

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

type CustomDomainResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p CustomDomainResourceOperationPredicate) Matches(input CustomDomainResource) bool {

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

type CustomizedAcceleratorResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p CustomizedAcceleratorResourceOperationPredicate) Matches(input CustomizedAcceleratorResource) bool {

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

type DeploymentResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p DeploymentResourceOperationPredicate) Matches(input DeploymentResource) bool {

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

type DevToolPortalResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p DevToolPortalResourceOperationPredicate) Matches(input DevToolPortalResource) bool {

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

type EurekaServerResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p EurekaServerResourceOperationPredicate) Matches(input EurekaServerResource) bool {

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

type GatewayCustomDomainResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p GatewayCustomDomainResourceOperationPredicate) Matches(input GatewayCustomDomainResource) bool {

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

type GatewayResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p GatewayResourceOperationPredicate) Matches(input GatewayResource) bool {

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

type GatewayRouteConfigResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p GatewayRouteConfigResourceOperationPredicate) Matches(input GatewayRouteConfigResource) bool {

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

type PredefinedAcceleratorResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p PredefinedAcceleratorResourceOperationPredicate) Matches(input PredefinedAcceleratorResource) bool {

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

type ResourceSkuOperationPredicate struct {
	Name         *string
	ResourceType *string
	Tier         *string
}

func (p ResourceSkuOperationPredicate) Matches(input ResourceSku) bool {

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.ResourceType != nil && (input.ResourceType == nil || *p.ResourceType != *input.ResourceType) {
		return false
	}

	if p.Tier != nil && (input.Tier == nil || *p.Tier != *input.Tier) {
		return false
	}

	return true
}

type ServiceRegistryResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ServiceRegistryResourceOperationPredicate) Matches(input ServiceRegistryResource) bool {

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

type ServiceResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ServiceResourceOperationPredicate) Matches(input ServiceResource) bool {

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

type StorageResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p StorageResourceOperationPredicate) Matches(input StorageResource) bool {

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

type SupportedApmTypeOperationPredicate struct {
	Name *string
}

func (p SupportedApmTypeOperationPredicate) Matches(input SupportedApmType) bool {

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	return true
}

type SupportedBuildpackResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p SupportedBuildpackResourceOperationPredicate) Matches(input SupportedBuildpackResource) bool {

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

type SupportedServerVersionOperationPredicate struct {
	Server  *string
	Value   *string
	Version *string
}

func (p SupportedServerVersionOperationPredicate) Matches(input SupportedServerVersion) bool {

	if p.Server != nil && (input.Server == nil || *p.Server != *input.Server) {
		return false
	}

	if p.Value != nil && (input.Value == nil || *p.Value != *input.Value) {
		return false
	}

	if p.Version != nil && (input.Version == nil || *p.Version != *input.Version) {
		return false
	}

	return true
}

type SupportedStackResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p SupportedStackResourceOperationPredicate) Matches(input SupportedStackResource) bool {

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
