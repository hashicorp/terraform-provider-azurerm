package contact

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactsProperties struct {
	AntennaConfiguration    *ContactsPropertiesAntennaConfiguration `json:"antennaConfiguration,omitempty"`
	ContactProfile          ResourceReference                       `json:"contactProfile"`
	EndAzimuthDegrees       *float64                                `json:"endAzimuthDegrees,omitempty"`
	EndElevationDegrees     *float64                                `json:"endElevationDegrees,omitempty"`
	ErrorMessage            *string                                 `json:"errorMessage,omitempty"`
	GroundStationName       string                                  `json:"groundStationName"`
	MaximumElevationDegrees *float64                                `json:"maximumElevationDegrees,omitempty"`
	ProvisioningState       *ProvisioningState                      `json:"provisioningState,omitempty"`
	ReservationEndTime      string                                  `json:"reservationEndTime"`
	ReservationStartTime    string                                  `json:"reservationStartTime"`
	RxEndTime               *string                                 `json:"rxEndTime,omitempty"`
	RxStartTime             *string                                 `json:"rxStartTime,omitempty"`
	StartAzimuthDegrees     *float64                                `json:"startAzimuthDegrees,omitempty"`
	StartElevationDegrees   *float64                                `json:"startElevationDegrees,omitempty"`
	Status                  *ContactsStatus                         `json:"status,omitempty"`
	TxEndTime               *string                                 `json:"txEndTime,omitempty"`
	TxStartTime             *string                                 `json:"txStartTime,omitempty"`
}

func (o *ContactsProperties) GetReservationEndTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.ReservationEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactsProperties) SetReservationEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ReservationEndTime = formatted
}

func (o *ContactsProperties) GetReservationStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.ReservationStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactsProperties) SetReservationStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ReservationStartTime = formatted
}

func (o *ContactsProperties) GetRxEndTimeAsTime() (*time.Time, error) {
	if o.RxEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RxEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactsProperties) SetRxEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RxEndTime = &formatted
}

func (o *ContactsProperties) GetRxStartTimeAsTime() (*time.Time, error) {
	if o.RxStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RxStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactsProperties) SetRxStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RxStartTime = &formatted
}

func (o *ContactsProperties) GetTxEndTimeAsTime() (*time.Time, error) {
	if o.TxEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TxEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactsProperties) SetTxEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TxEndTime = &formatted
}

func (o *ContactsProperties) GetTxStartTimeAsTime() (*time.Time, error) {
	if o.TxStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TxStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContactsProperties) SetTxStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TxStartTime = &formatted
}
