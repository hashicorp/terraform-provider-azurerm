package contact

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactInstanceProperties struct {
	EndAzimuthDegrees       *float64 `json:"endAzimuthDegrees,omitempty"`
	EndElevationDegrees     *float64 `json:"endElevationDegrees,omitempty"`
	MaximumElevationDegrees *float64 `json:"maximumElevationDegrees,omitempty"`
	RxEndTime               *string  `json:"rxEndTime,omitempty"`
	RxStartTime             *string  `json:"rxStartTime,omitempty"`
	StartAzimuthDegrees     *float64 `json:"startAzimuthDegrees,omitempty"`
	StartElevationDegrees   *float64 `json:"startElevationDegrees,omitempty"`
	TxEndTime               *string  `json:"txEndTime,omitempty"`
	TxStartTime             *string  `json:"txStartTime,omitempty"`
}

func (o *ContactInstanceProperties) GetRxEndTimeAsTime() (*time.Time, error) {
	if o.RxEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RxEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactInstanceProperties) SetRxEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RxEndTime = &formatted
}

func (o *ContactInstanceProperties) GetRxStartTimeAsTime() (*time.Time, error) {
	if o.RxStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RxStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactInstanceProperties) SetRxStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RxStartTime = &formatted
}

func (o *ContactInstanceProperties) GetTxEndTimeAsTime() (*time.Time, error) {
	if o.TxEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TxEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactInstanceProperties) SetTxEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TxEndTime = &formatted
}

func (o *ContactInstanceProperties) GetTxStartTimeAsTime() (*time.Time, error) {
	if o.TxStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TxStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactInstanceProperties) SetTxStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TxStartTime = &formatted
}
