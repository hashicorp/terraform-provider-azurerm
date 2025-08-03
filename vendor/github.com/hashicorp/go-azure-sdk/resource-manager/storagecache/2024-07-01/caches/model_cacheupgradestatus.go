package caches

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheUpgradeStatus struct {
	CurrentFirmwareVersion *string             `json:"currentFirmwareVersion,omitempty"`
	FirmwareUpdateDeadline *string             `json:"firmwareUpdateDeadline,omitempty"`
	FirmwareUpdateStatus   *FirmwareStatusType `json:"firmwareUpdateStatus,omitempty"`
	LastFirmwareUpdate     *string             `json:"lastFirmwareUpdate,omitempty"`
	PendingFirmwareVersion *string             `json:"pendingFirmwareVersion,omitempty"`
}

func (o *CacheUpgradeStatus) GetFirmwareUpdateDeadlineAsTime() (*time.Time, error) {
	if o.FirmwareUpdateDeadline == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FirmwareUpdateDeadline, "2006-01-02T15:04:05Z07:00")
}

func (o *CacheUpgradeStatus) SetFirmwareUpdateDeadlineAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FirmwareUpdateDeadline = &formatted
}

func (o *CacheUpgradeStatus) GetLastFirmwareUpdateAsTime() (*time.Time, error) {
	if o.LastFirmwareUpdate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastFirmwareUpdate, "2006-01-02T15:04:05Z07:00")
}

func (o *CacheUpgradeStatus) SetLastFirmwareUpdateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastFirmwareUpdate = &formatted
}
