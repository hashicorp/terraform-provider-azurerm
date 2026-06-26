package nodereports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscNodeReportOperationPredicate struct {
	ConfigurationVersion *string
	EndTime              *string
	HostName             *string
	Id                   *string
	LastModifiedTime     *string
	NumberOfResources    *int64
	RawErrors            *string
	RebootRequested      *string
	RefreshMode          *string
	ReportFormatVersion  *string
	ReportId             *string
	StartTime            *string
	Status               *string
	Type                 *string
}

func (p DscNodeReportOperationPredicate) Matches(input DscNodeReport) bool {

	if p.ConfigurationVersion != nil && (input.ConfigurationVersion == nil || *p.ConfigurationVersion != *input.ConfigurationVersion) {
		return false
	}

	if p.EndTime != nil && (input.EndTime == nil || *p.EndTime != *input.EndTime) {
		return false
	}

	if p.HostName != nil && (input.HostName == nil || *p.HostName != *input.HostName) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.LastModifiedTime != nil && (input.LastModifiedTime == nil || *p.LastModifiedTime != *input.LastModifiedTime) {
		return false
	}

	if p.NumberOfResources != nil && (input.NumberOfResources == nil || *p.NumberOfResources != *input.NumberOfResources) {
		return false
	}

	if p.RawErrors != nil && (input.RawErrors == nil || *p.RawErrors != *input.RawErrors) {
		return false
	}

	if p.RebootRequested != nil && (input.RebootRequested == nil || *p.RebootRequested != *input.RebootRequested) {
		return false
	}

	if p.RefreshMode != nil && (input.RefreshMode == nil || *p.RefreshMode != *input.RefreshMode) {
		return false
	}

	if p.ReportFormatVersion != nil && (input.ReportFormatVersion == nil || *p.ReportFormatVersion != *input.ReportFormatVersion) {
		return false
	}

	if p.ReportId != nil && (input.ReportId == nil || *p.ReportId != *input.ReportId) {
		return false
	}

	if p.StartTime != nil && (input.StartTime == nil || *p.StartTime != *input.StartTime) {
		return false
	}

	if p.Status != nil && (input.Status == nil || *p.Status != *input.Status) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}
