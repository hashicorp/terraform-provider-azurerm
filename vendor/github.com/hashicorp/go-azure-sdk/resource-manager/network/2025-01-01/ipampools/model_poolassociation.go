package ipampools

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolAssociation struct {
	AddressPrefixes             *[]string `json:"addressPrefixes,omitempty"`
	CreatedAt                   *string   `json:"createdAt,omitempty"`
	Description                 *string   `json:"description,omitempty"`
	NumberOfReservedIPAddresses *string   `json:"numberOfReservedIPAddresses,omitempty"`
	PoolId                      *string   `json:"poolId,omitempty"`
	ReservationExpiresAt        *string   `json:"reservationExpiresAt,omitempty"`
	ReservedPrefixes            *[]string `json:"reservedPrefixes,omitempty"`
	ResourceId                  string    `json:"resourceId"`
	TotalNumberOfIPAddresses    *string   `json:"totalNumberOfIPAddresses,omitempty"`
}

func (o *PoolAssociation) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *PoolAssociation) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *PoolAssociation) GetReservationExpiresAtAsTime() (*time.Time, error) {
	if o.ReservationExpiresAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ReservationExpiresAt, "2006-01-02T15:04:05Z07:00")
}

func (o *PoolAssociation) SetReservationExpiresAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ReservationExpiresAt = &formatted
}
