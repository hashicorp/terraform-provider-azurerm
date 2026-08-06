package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedPartnerResourcesListFormatOperationPredicate struct {
}

func (p ConnectedPartnerResourcesListFormatOperationPredicate) Matches(input ConnectedPartnerResourcesListFormat) bool {

	return true
}

type ElasticMonitorResourceOperationPredicate struct {
	Id       *string
	Kind     *string
	Location *string
	Name     *string
	Type     *string
}

func (p ElasticMonitorResourceOperationPredicate) Matches(input ElasticMonitorResource) bool {

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

type MonitoredResourceOperationPredicate struct {
	Id                  *string
	ReasonForLogsStatus *string
}

func (p MonitoredResourceOperationPredicate) Matches(input MonitoredResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.ReasonForLogsStatus != nil && (input.ReasonForLogsStatus == nil || *p.ReasonForLogsStatus != *input.ReasonForLogsStatus) {
		return false
	}

	return true
}

type VMResourcesOperationPredicate struct {
	VMResourceId *string
}

func (p VMResourcesOperationPredicate) Matches(input VMResources) bool {

	if p.VMResourceId != nil && (input.VMResourceId == nil || *p.VMResourceId != *input.VMResourceId) {
		return false
	}

	return true
}
