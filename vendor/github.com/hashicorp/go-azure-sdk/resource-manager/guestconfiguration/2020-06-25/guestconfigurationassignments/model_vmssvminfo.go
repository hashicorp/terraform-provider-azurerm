package guestconfigurationassignments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMSSVMInfo struct {
	ComplianceStatus      *ComplianceStatus `json:"complianceStatus,omitempty"`
	LastComplianceChecked *string           `json:"lastComplianceChecked,omitempty"`
	LatestReportId        *string           `json:"latestReportId,omitempty"`
	VMId                  *string           `json:"vmId,omitempty"`
	VMResourceId          *string           `json:"vmResourceId,omitempty"`
}

func (o *VMSSVMInfo) GetLastComplianceCheckedAsTime() (*time.Time, error) {
	if o.LastComplianceChecked == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastComplianceChecked, "2006-01-02T15:04:05Z07:00")
}

func (o *VMSSVMInfo) SetLastComplianceCheckedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastComplianceChecked = &formatted
}
