package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiKVReferenceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p ApiKVReferenceOperationPredicate) Matches(input ApiKVReference) bool {

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

type BackupItemOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p BackupItemOperationPredicate) Matches(input BackupItem) bool {

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

type ContinuousWebJobOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p ContinuousWebJobOperationPredicate) Matches(input ContinuousWebJob) bool {

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

type CsmDeploymentStatusOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p CsmDeploymentStatusOperationPredicate) Matches(input CsmDeploymentStatus) bool {

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

type CsmPublishingCredentialsPoliciesEntityOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p CsmPublishingCredentialsPoliciesEntityOperationPredicate) Matches(input CsmPublishingCredentialsPoliciesEntity) bool {

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

type DeploymentOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p DeploymentOperationPredicate) Matches(input Deployment) bool {

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

type FunctionEnvelopeOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p FunctionEnvelopeOperationPredicate) Matches(input FunctionEnvelope) bool {

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

type HostNameBindingOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p HostNameBindingOperationPredicate) Matches(input HostNameBinding) bool {

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

type IdentifierOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p IdentifierOperationPredicate) Matches(input Identifier) bool {

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

type PerfMonResponseOperationPredicate struct {
	Code    *string
	Message *string
}

func (p PerfMonResponseOperationPredicate) Matches(input PerfMonResponse) bool {

	if p.Code != nil && (input.Code == nil || *p.Code != *input.Code) {
		return false
	}

	if p.Message != nil && (input.Message == nil || *p.Message != *input.Message) {
		return false
	}

	return true
}

type ProcessInfoOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p ProcessInfoOperationPredicate) Matches(input ProcessInfo) bool {

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

type ProcessModuleInfoOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p ProcessModuleInfoOperationPredicate) Matches(input ProcessModuleInfo) bool {

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

type ProcessThreadInfoOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p ProcessThreadInfoOperationPredicate) Matches(input ProcessThreadInfo) bool {

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

type PublicCertificateOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p PublicCertificateOperationPredicate) Matches(input PublicCertificate) bool {

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

type SiteConfigResourceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p SiteConfigResourceOperationPredicate) Matches(input SiteConfigResource) bool {

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

type SiteConfigurationSnapshotInfoOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p SiteConfigurationSnapshotInfoOperationPredicate) Matches(input SiteConfigurationSnapshotInfo) bool {

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

type SiteContainerOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p SiteContainerOperationPredicate) Matches(input SiteContainer) bool {

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

type SiteExtensionInfoOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p SiteExtensionInfoOperationPredicate) Matches(input SiteExtensionInfo) bool {

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

type SlotDifferenceOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p SlotDifferenceOperationPredicate) Matches(input SlotDifference) bool {

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

type SnapshotOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p SnapshotOperationPredicate) Matches(input Snapshot) bool {

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

type TriggeredJobHistoryOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p TriggeredJobHistoryOperationPredicate) Matches(input TriggeredJobHistory) bool {

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

type TriggeredWebJobOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p TriggeredWebJobOperationPredicate) Matches(input TriggeredWebJob) bool {

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

type WebJobOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p WebJobOperationPredicate) Matches(input WebJob) bool {

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

type WebSiteInstanceStatusOperationPredicate struct {
	Id   *string
	Kind *string
	Name *string
	Type *string
}

func (p WebSiteInstanceStatusOperationPredicate) Matches(input WebSiteInstanceStatus) bool {

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

type WorkflowEnvelopeOperationPredicate struct {
	Id       *string
	Kind     *string
	Location *string
	Name     *string
	Type     *string
}

func (p WorkflowEnvelopeOperationPredicate) Matches(input WorkflowEnvelope) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Kind != nil && (input.Kind == nil || *p.Kind != *input.Kind) {
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
