package restorables

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableDatabaseAccountProperties struct {
	AccountName          *string                       `json:"accountName,omitempty"`
	ApiType              *ApiType                      `json:"apiType,omitempty"`
	CreationTime         *string                       `json:"creationTime,omitempty"`
	DeletionTime         *string                       `json:"deletionTime,omitempty"`
	OldestRestorableTime *string                       `json:"oldestRestorableTime,omitempty"`
	RestorableLocations  *[]RestorableLocationResource `json:"restorableLocations,omitempty"`
}

func (o *RestorableDatabaseAccountProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorableDatabaseAccountProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *RestorableDatabaseAccountProperties) GetDeletionTimeAsTime() (*time.Time, error) {
	if o.DeletionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeletionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorableDatabaseAccountProperties) SetDeletionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeletionTime = &formatted
}

func (o *RestorableDatabaseAccountProperties) GetOldestRestorableTimeAsTime() (*time.Time, error) {
	if o.OldestRestorableTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OldestRestorableTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorableDatabaseAccountProperties) SetOldestRestorableTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OldestRestorableTime = &formatted
}
