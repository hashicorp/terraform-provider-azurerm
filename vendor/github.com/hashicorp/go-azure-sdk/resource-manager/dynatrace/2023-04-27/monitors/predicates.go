package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceInfoOperationPredicate struct {
	HostGroup  *string
	HostName   *string
	ResourceId *string
	Version    *string
}

func (p AppServiceInfoOperationPredicate) Matches(input AppServiceInfo) bool {

	if p.HostGroup != nil && (input.HostGroup == nil || *p.HostGroup != *input.HostGroup) {
		return false
	}

	if p.HostName != nil && (input.HostName == nil || *p.HostName != *input.HostName) {
		return false
	}

	if p.ResourceId != nil && (input.ResourceId == nil || *p.ResourceId != *input.ResourceId) {
		return false
	}

	if p.Version != nil && (input.Version == nil || *p.Version != *input.Version) {
		return false
	}

	return true
}

type LinkableEnvironmentResponseOperationPredicate struct {
	EnvironmentId   *string
	EnvironmentName *string
}

func (p LinkableEnvironmentResponseOperationPredicate) Matches(input LinkableEnvironmentResponse) bool {

	if p.EnvironmentId != nil && (input.EnvironmentId == nil || *p.EnvironmentId != *input.EnvironmentId) {
		return false
	}

	if p.EnvironmentName != nil && (input.EnvironmentName == nil || *p.EnvironmentName != *input.EnvironmentName) {
		return false
	}

	return true
}

type MonitorResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p MonitorResourceOperationPredicate) Matches(input MonitorResource) bool {

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

type MonitoredResourceOperationPredicate struct {
	Id                     *string
	ReasonForLogsStatus    *string
	ReasonForMetricsStatus *string
}

func (p MonitoredResourceOperationPredicate) Matches(input MonitoredResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.ReasonForLogsStatus != nil && (input.ReasonForLogsStatus == nil || *p.ReasonForLogsStatus != *input.ReasonForLogsStatus) {
		return false
	}

	if p.ReasonForMetricsStatus != nil && (input.ReasonForMetricsStatus == nil || *p.ReasonForMetricsStatus != *input.ReasonForMetricsStatus) {
		return false
	}

	return true
}

type VMInfoOperationPredicate struct {
	HostGroup  *string
	HostName   *string
	ResourceId *string
	Version    *string
}

func (p VMInfoOperationPredicate) Matches(input VMInfo) bool {

	if p.HostGroup != nil && (input.HostGroup == nil || *p.HostGroup != *input.HostGroup) {
		return false
	}

	if p.HostName != nil && (input.HostName == nil || *p.HostName != *input.HostName) {
		return false
	}

	if p.ResourceId != nil && (input.ResourceId == nil || *p.ResourceId != *input.ResourceId) {
		return false
	}

	if p.Version != nil && (input.Version == nil || *p.Version != *input.Version) {
		return false
	}

	return true
}
