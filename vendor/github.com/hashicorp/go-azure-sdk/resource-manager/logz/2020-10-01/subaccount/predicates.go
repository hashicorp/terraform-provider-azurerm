package subaccount

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogzMonitorResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p LogzMonitorResourceOperationPredicate) Matches(input LogzMonitorResource) bool {

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
	SendingLogs            *bool
	SendingMetrics         *bool
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

	if p.SendingLogs != nil && (input.SendingLogs == nil || *p.SendingLogs != *input.SendingLogs) {
		return false
	}

	if p.SendingMetrics != nil && (input.SendingMetrics == nil || *p.SendingMetrics != *input.SendingMetrics) {
		return false
	}

	return true
}
