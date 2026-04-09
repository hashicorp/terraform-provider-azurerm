package catalogs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CatalogProperties struct {
	AdoGit             *GitCatalog             `json:"adoGit,omitempty"`
	ConnectionState    *CatalogConnectionState `json:"connectionState,omitempty"`
	GitHub             *GitCatalog             `json:"gitHub,omitempty"`
	LastConnectionTime *string                 `json:"lastConnectionTime,omitempty"`
	LastSyncStats      *SyncStats              `json:"lastSyncStats,omitempty"`
	LastSyncTime       *string                 `json:"lastSyncTime,omitempty"`
	ProvisioningState  *ProvisioningState      `json:"provisioningState,omitempty"`
	SyncState          *CatalogSyncState       `json:"syncState,omitempty"`
	SyncType           *CatalogSyncType        `json:"syncType,omitempty"`
	Tags               *map[string]string      `json:"tags,omitempty"`
}

func (o *CatalogProperties) GetLastConnectionTimeAsTime() (*time.Time, error) {
	if o.LastConnectionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastConnectionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CatalogProperties) SetLastConnectionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastConnectionTime = &formatted
}

func (o *CatalogProperties) GetLastSyncTimeAsTime() (*time.Time, error) {
	if o.LastSyncTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSyncTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CatalogProperties) SetLastSyncTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSyncTime = &formatted
}
