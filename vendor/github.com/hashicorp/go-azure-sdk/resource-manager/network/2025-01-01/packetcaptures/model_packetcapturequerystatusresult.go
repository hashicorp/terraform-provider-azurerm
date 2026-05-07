package packetcaptures

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCaptureQueryStatusResult struct {
	CaptureStartTime    *string    `json:"captureStartTime,omitempty"`
	Id                  *string    `json:"id,omitempty"`
	Name                *string    `json:"name,omitempty"`
	PacketCaptureError  *[]PcError `json:"packetCaptureError,omitempty"`
	PacketCaptureStatus *PcStatus  `json:"packetCaptureStatus,omitempty"`
	StopReason          *string    `json:"stopReason,omitempty"`
}

func (o *PacketCaptureQueryStatusResult) GetCaptureStartTimeAsTime() (*time.Time, error) {
	if o.CaptureStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CaptureStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PacketCaptureQueryStatusResult) SetCaptureStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CaptureStartTime = &formatted
}
