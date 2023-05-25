package agreements

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OldAgreementProperties struct {
	CancelDate *string `json:"cancelDate,omitempty"`
	Id         *string `json:"id,omitempty"`
	Offer      *string `json:"offer,omitempty"`
	Publisher  *string `json:"publisher,omitempty"`
	SignDate   *string `json:"signDate,omitempty"`
	State      *State  `json:"state,omitempty"`
}

func (o *OldAgreementProperties) GetCancelDateAsTime() (*time.Time, error) {
	if o.CancelDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CancelDate, "2006-01-02T15:04:05Z07:00")
}

func (o *OldAgreementProperties) SetCancelDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CancelDate = &formatted
}

func (o *OldAgreementProperties) GetSignDateAsTime() (*time.Time, error) {
	if o.SignDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SignDate, "2006-01-02T15:04:05Z07:00")
}

func (o *OldAgreementProperties) SetSignDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SignDate = &formatted
}
