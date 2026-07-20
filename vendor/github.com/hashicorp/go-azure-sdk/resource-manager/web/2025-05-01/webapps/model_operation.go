package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Operation struct {
	CreatedTime          *string          `json:"createdTime,omitempty"`
	Errors               *[]ErrorEntity   `json:"errors,omitempty"`
	ExpirationTime       *string          `json:"expirationTime,omitempty"`
	GeoMasterOperationId *string          `json:"geoMasterOperationId,omitempty"`
	Id                   *string          `json:"id,omitempty"`
	ModifiedTime         *string          `json:"modifiedTime,omitempty"`
	Name                 *string          `json:"name,omitempty"`
	Status               *OperationStatus `json:"status,omitempty"`
}

func (o *Operation) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Operation) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}

func (o *Operation) GetExpirationTimeAsTime() (*time.Time, error) {
	if o.ExpirationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Operation) SetExpirationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTime = &formatted
}

func (o *Operation) GetModifiedTimeAsTime() (*time.Time, error) {
	if o.ModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Operation) SetModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ModifiedTime = &formatted
}
