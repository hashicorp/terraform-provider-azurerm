package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServicesListResponseOperationPredicate struct {
	NextLink *string
}

func (p AppServicesListResponseOperationPredicate) Matches(input AppServicesListResponse) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}

type MonitoredResourceListResponseOperationPredicate struct {
	NextLink *string
}

func (p MonitoredResourceListResponseOperationPredicate) Matches(input MonitoredResourceListResponse) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}

type NewRelicMonitorResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p NewRelicMonitorResourceOperationPredicate) Matches(input NewRelicMonitorResource) bool {

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

type VMHostsListResponseOperationPredicate struct {
	NextLink *string
}

func (p VMHostsListResponseOperationPredicate) Matches(input VMHostsListResponse) bool {

	if p.NextLink != nil && (input.NextLink == nil || *p.NextLink != *input.NextLink) {
		return false
	}

	return true
}
