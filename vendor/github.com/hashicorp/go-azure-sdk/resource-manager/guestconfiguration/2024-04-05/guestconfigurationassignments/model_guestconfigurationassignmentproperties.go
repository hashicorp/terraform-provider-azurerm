package guestconfigurationassignments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GuestConfigurationAssignmentProperties struct {
	AssignmentHash              *string                       `json:"assignmentHash,omitempty"`
	ComplianceStatus            *ComplianceStatus             `json:"complianceStatus,omitempty"`
	Context                     *string                       `json:"context,omitempty"`
	GuestConfiguration          *GuestConfigurationNavigation `json:"guestConfiguration,omitempty"`
	LastComplianceStatusChecked *string                       `json:"lastComplianceStatusChecked,omitempty"`
	LatestAssignmentReport      *AssignmentReport             `json:"latestAssignmentReport,omitempty"`
	LatestReportId              *string                       `json:"latestReportId,omitempty"`
	ParameterHash               *string                       `json:"parameterHash,omitempty"`
	ProvisioningState           *ProvisioningState            `json:"provisioningState,omitempty"`
	ResourceType                *string                       `json:"resourceType,omitempty"`
	TargetResourceId            *string                       `json:"targetResourceId,omitempty"`
	VMSSVMList                  *[]VMSSVMInfo                 `json:"vmssVMList,omitempty"`
}

func (o *GuestConfigurationAssignmentProperties) GetLastComplianceStatusCheckedAsTime() (*time.Time, error) {
	if o.LastComplianceStatusChecked == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastComplianceStatusChecked, "2006-01-02T15:04:05Z07:00")
}

func (o *GuestConfigurationAssignmentProperties) SetLastComplianceStatusCheckedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastComplianceStatusChecked = &formatted
}
