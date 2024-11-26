package connectionmonitors

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionStateSnapshot struct {
	AvgLatencyInMs  *int64             `json:"avgLatencyInMs,omitempty"`
	ConnectionState *ConnectionState   `json:"connectionState,omitempty"`
	EndTime         *string            `json:"endTime,omitempty"`
	EvaluationState *EvaluationState   `json:"evaluationState,omitempty"`
	Hops            *[]ConnectivityHop `json:"hops,omitempty"`
	MaxLatencyInMs  *int64             `json:"maxLatencyInMs,omitempty"`
	MinLatencyInMs  *int64             `json:"minLatencyInMs,omitempty"`
	ProbesFailed    *int64             `json:"probesFailed,omitempty"`
	ProbesSent      *int64             `json:"probesSent,omitempty"`
	StartTime       *string            `json:"startTime,omitempty"`
}

func (o *ConnectionStateSnapshot) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectionStateSnapshot) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *ConnectionStateSnapshot) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectionStateSnapshot) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
