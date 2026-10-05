package sqlvirtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TroubleshootingStatus struct {
	EndTimeUtc              *string                              `json:"endTimeUtc,omitempty"`
	LastTriggerTimeUtc      *string                              `json:"lastTriggerTimeUtc,omitempty"`
	Properties              *TroubleshootingAdditionalProperties `json:"properties,omitempty"`
	RootCause               *string                              `json:"rootCause,omitempty"`
	StartTimeUtc            *string                              `json:"startTimeUtc,omitempty"`
	TroubleshootingScenario *TroubleshootingScenario             `json:"troubleshootingScenario,omitempty"`
}

func (o *TroubleshootingStatus) GetEndTimeUtcAsTime() (*time.Time, error) {
	if o.EndTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *TroubleshootingStatus) SetEndTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTimeUtc = &formatted
}

func (o *TroubleshootingStatus) GetLastTriggerTimeUtcAsTime() (*time.Time, error) {
	if o.LastTriggerTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastTriggerTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *TroubleshootingStatus) SetLastTriggerTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastTriggerTimeUtc = &formatted
}

func (o *TroubleshootingStatus) GetStartTimeUtcAsTime() (*time.Time, error) {
	if o.StartTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *TroubleshootingStatus) SetStartTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTimeUtc = &formatted
}
