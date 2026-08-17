package restorables

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableLocationResource struct {
	CreationTime                      *string `json:"creationTime,omitempty"`
	DeletionTime                      *string `json:"deletionTime,omitempty"`
	LocationName                      *string `json:"locationName,omitempty"`
	RegionalDatabaseAccountInstanceId *string `json:"regionalDatabaseAccountInstanceId,omitempty"`
}

func (o *RestorableLocationResource) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorableLocationResource) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *RestorableLocationResource) GetDeletionTimeAsTime() (*time.Time, error) {
	if o.DeletionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeletionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorableLocationResource) SetDeletionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeletionTime = &formatted
}
