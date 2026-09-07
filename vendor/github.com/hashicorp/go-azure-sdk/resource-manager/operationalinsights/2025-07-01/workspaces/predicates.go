package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementGroupOperationPredicate struct {
}

func (p ManagementGroupOperationPredicate) Matches(input ManagementGroup) bool {

	return true
}

type UsageMetricOperationPredicate struct {
	CurrentValue  *float64
	Limit         *float64
	NextResetTime *string
	QuotaPeriod   *string
	Unit          *string
}

func (p UsageMetricOperationPredicate) Matches(input UsageMetric) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil || *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	if p.NextResetTime != nil && (input.NextResetTime == nil || *p.NextResetTime != *input.NextResetTime) {
		return false
	}

	if p.QuotaPeriod != nil && (input.QuotaPeriod == nil || *p.QuotaPeriod != *input.QuotaPeriod) {
		return false
	}

	if p.Unit != nil && (input.Unit == nil || *p.Unit != *input.Unit) {
		return false
	}

	return true
}

type WorkspaceOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p WorkspaceOperationPredicate) Matches(input Workspace) bool {

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
