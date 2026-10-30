package protectioncontainers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MabContainerExtendedInfo struct {
	BackupItemType   *BackupItemType `json:"backupItemType,omitempty"`
	BackupItems      *[]string       `json:"backupItems,omitempty"`
	LastBackupStatus *string         `json:"lastBackupStatus,omitempty"`
	LastRefreshedAt  *string         `json:"lastRefreshedAt,omitempty"`
	PolicyName       *string         `json:"policyName,omitempty"`
}

func (o *MabContainerExtendedInfo) GetLastRefreshedAtAsTime() (*time.Time, error) {
	if o.LastRefreshedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRefreshedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *MabContainerExtendedInfo) SetLastRefreshedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRefreshedAt = &formatted
}
