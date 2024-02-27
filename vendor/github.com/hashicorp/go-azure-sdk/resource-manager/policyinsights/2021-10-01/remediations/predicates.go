package remediations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RemediationOperationPredicate) Matches(input Remediation) bool {

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

type RemediationDeploymentOperationPredicate struct {
	CreatedOn            *string
	DeploymentId         *string
	LastUpdatedOn        *string
	RemediatedResourceId *string
	ResourceLocation     *string
	Status               *string
}

func (p RemediationDeploymentOperationPredicate) Matches(input RemediationDeployment) bool {

	if p.CreatedOn != nil && (input.CreatedOn == nil || *p.CreatedOn != *input.CreatedOn) {
		return false
	}

	if p.DeploymentId != nil && (input.DeploymentId == nil || *p.DeploymentId != *input.DeploymentId) {
		return false
	}

	if p.LastUpdatedOn != nil && (input.LastUpdatedOn == nil || *p.LastUpdatedOn != *input.LastUpdatedOn) {
		return false
	}

	if p.RemediatedResourceId != nil && (input.RemediatedResourceId == nil || *p.RemediatedResourceId != *input.RemediatedResourceId) {
		return false
	}

	if p.ResourceLocation != nil && (input.ResourceLocation == nil || *p.ResourceLocation != *input.ResourceLocation) {
		return false
	}

	if p.Status != nil && (input.Status == nil || *p.Status != *input.Status) {
		return false
	}

	return true
}
