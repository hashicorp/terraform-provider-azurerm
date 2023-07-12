package watchlistitems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatchlistItemProperties struct {
	Created           *string      `json:"created,omitempty"`
	CreatedBy         *UserInfo    `json:"createdBy,omitempty"`
	EntityMapping     *interface{} `json:"entityMapping,omitempty"`
	IsDeleted         *bool        `json:"isDeleted,omitempty"`
	ItemsKeyValue     interface{}  `json:"itemsKeyValue"`
	TenantId          *string      `json:"tenantId,omitempty"`
	Updated           *string      `json:"updated,omitempty"`
	UpdatedBy         *UserInfo    `json:"updatedBy,omitempty"`
	WatchlistItemId   *string      `json:"watchlistItemId,omitempty"`
	WatchlistItemType *string      `json:"watchlistItemType,omitempty"`
}

func (o *WatchlistItemProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *WatchlistItemProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *WatchlistItemProperties) GetUpdatedAsTime() (*time.Time, error) {
	if o.Updated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Updated, "2006-01-02T15:04:05Z07:00")
}

func (o *WatchlistItemProperties) SetUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Updated = &formatted
}
