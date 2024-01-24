package catalogs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CatalogProperties struct {
	AdoGit            *GitCatalog        `json:"adoGit,omitempty"`
	GitHub            *GitCatalog        `json:"gitHub,omitempty"`
	LastSyncTime      *string            `json:"lastSyncTime,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	SyncState         *CatalogSyncState  `json:"syncState,omitempty"`
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
