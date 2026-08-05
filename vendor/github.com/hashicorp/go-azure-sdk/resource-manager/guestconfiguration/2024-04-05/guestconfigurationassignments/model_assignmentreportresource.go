package guestconfigurationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssignmentReportResource struct {
	ComplianceStatus *ComplianceStatus                           `json:"complianceStatus,omitempty"`
	Properties       *interface{}                                `json:"properties,omitempty"`
	Reasons          *[]AssignmentReportResourceComplianceReason `json:"reasons,omitempty"`
	ResourceId       *string                                     `json:"resourceId,omitempty"`
}
