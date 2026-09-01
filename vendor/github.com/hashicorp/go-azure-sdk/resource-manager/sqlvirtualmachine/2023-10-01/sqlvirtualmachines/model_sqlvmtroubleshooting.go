package sqlvirtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlVMTroubleshooting struct {
	EndTimeUtc               *string                              `json:"endTimeUtc,omitempty"`
	Properties               *TroubleshootingAdditionalProperties `json:"properties,omitempty"`
	StartTimeUtc             *string                              `json:"startTimeUtc,omitempty"`
	TroubleshootingScenario  *TroubleshootingScenario             `json:"troubleshootingScenario,omitempty"`
	VirtualMachineResourceId *string                              `json:"virtualMachineResourceId,omitempty"`
}

func (o *SqlVMTroubleshooting) GetEndTimeUtcAsTime() (*time.Time, error) {
	if o.EndTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *SqlVMTroubleshooting) SetEndTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTimeUtc = &formatted
}

func (o *SqlVMTroubleshooting) GetStartTimeUtcAsTime() (*time.Time, error) {
	if o.StartTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *SqlVMTroubleshooting) SetStartTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTimeUtc = &formatted
}
