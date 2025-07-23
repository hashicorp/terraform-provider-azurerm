package storagetaskassignments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggerParametersUpdate struct {
	EndBy        *string       `json:"endBy,omitempty"`
	Interval     *int64        `json:"interval,omitempty"`
	IntervalUnit *IntervalUnit `json:"intervalUnit,omitempty"`
	StartFrom    *string       `json:"startFrom,omitempty"`
	StartOn      *string       `json:"startOn,omitempty"`
}

func (o *TriggerParametersUpdate) GetEndByAsTime() (*time.Time, error) {
	if o.EndBy == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndBy, "2006-01-02T15:04:05Z07:00")
}

func (o *TriggerParametersUpdate) SetEndByAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndBy = &formatted
}

func (o *TriggerParametersUpdate) GetStartFromAsTime() (*time.Time, error) {
	if o.StartFrom == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartFrom, "2006-01-02T15:04:05Z07:00")
}

func (o *TriggerParametersUpdate) SetStartFromAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartFrom = &formatted
}

func (o *TriggerParametersUpdate) GetStartOnAsTime() (*time.Time, error) {
	if o.StartOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartOn, "2006-01-02T15:04:05Z07:00")
}

func (o *TriggerParametersUpdate) SetStartOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartOn = &formatted
}
