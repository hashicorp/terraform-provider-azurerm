package deletedconfigurationstores

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedConfigurationStoreProperties struct {
	ConfigurationStoreId   *string            `json:"configurationStoreId,omitempty"`
	DeletionDate           *string            `json:"deletionDate,omitempty"`
	Location               *string            `json:"location,omitempty"`
	PurgeProtectionEnabled *bool              `json:"purgeProtectionEnabled,omitempty"`
	ScheduledPurgeDate     *string            `json:"scheduledPurgeDate,omitempty"`
	Tags                   *map[string]string `json:"tags,omitempty"`
}

func (o *DeletedConfigurationStoreProperties) GetDeletionDateAsTime() (*time.Time, error) {
	if o.DeletionDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeletionDate, "2006-01-02T15:04:05Z07:00")
}

func (o *DeletedConfigurationStoreProperties) SetDeletionDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeletionDate = &formatted
}

func (o *DeletedConfigurationStoreProperties) GetScheduledPurgeDateAsTime() (*time.Time, error) {
	if o.ScheduledPurgeDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ScheduledPurgeDate, "2006-01-02T15:04:05Z07:00")
}

func (o *DeletedConfigurationStoreProperties) SetScheduledPurgeDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ScheduledPurgeDate = &formatted
}
