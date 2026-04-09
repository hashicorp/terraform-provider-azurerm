package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceInfoOperationPredicate struct {
	AgentStatus     *string
	AgentVersion    *string
	AzureResourceId *string
}

func (p AppServiceInfoOperationPredicate) Matches(input AppServiceInfo) bool {

	if p.AgentStatus != nil && (input.AgentStatus == nil || *p.AgentStatus != *input.AgentStatus) {
		return false
	}

	if p.AgentVersion != nil && (input.AgentVersion == nil || *p.AgentVersion != *input.AgentVersion) {
		return false
	}

	if p.AzureResourceId != nil && (input.AzureResourceId == nil || *p.AzureResourceId != *input.AzureResourceId) {
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

type VMInfoOperationPredicate struct {
	AgentStatus  *string
	AgentVersion *string
	VMId         *string
}

func (p VMInfoOperationPredicate) Matches(input VMInfo) bool {

	if p.AgentStatus != nil && (input.AgentStatus == nil || *p.AgentStatus != *input.AgentStatus) {
		return false
	}

	if p.AgentVersion != nil && (input.AgentVersion == nil || *p.AgentVersion != *input.AgentVersion) {
		return false
	}

	if p.VMId != nil && (input.VMId == nil || *p.VMId != *input.VMId) {
		return false
	}

	return true
}
