package recoverypoint

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPointDataStoreDetails struct {
	CreationTime          *string            `json:"creationTime,omitempty"`
	ExpiryTime            *string            `json:"expiryTime,omitempty"`
	Id                    *string            `json:"id,omitempty"`
	MetaData              *string            `json:"metaData,omitempty"`
	RehydrationExpiryTime *string            `json:"rehydrationExpiryTime,omitempty"`
	RehydrationStatus     *RehydrationStatus `json:"rehydrationStatus,omitempty"`
	State                 *string            `json:"state,omitempty"`
	Type                  *string            `json:"type,omitempty"`
	Visible               *bool              `json:"visible,omitempty"`
}

func (o *RecoveryPointDataStoreDetails) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RecoveryPointDataStoreDetails) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *RecoveryPointDataStoreDetails) GetExpiryTimeAsTime() (*time.Time, error) {
	if o.ExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RecoveryPointDataStoreDetails) SetExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryTime = &formatted
}

func (o *RecoveryPointDataStoreDetails) GetRehydrationExpiryTimeAsTime() (*time.Time, error) {
	if o.RehydrationExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RehydrationExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RecoveryPointDataStoreDetails) SetRehydrationExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RehydrationExpiryTime = &formatted
}
