package watchlists

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatchlistProperties struct {
	ContentType         *string   `json:"contentType,omitempty"`
	Created             *string   `json:"created,omitempty"`
	CreatedBy           *UserInfo `json:"createdBy,omitempty"`
	DefaultDuration     *string   `json:"defaultDuration,omitempty"`
	Description         *string   `json:"description,omitempty"`
	DisplayName         string    `json:"displayName"`
	IsDeleted           *bool     `json:"isDeleted,omitempty"`
	ItemsSearchKey      string    `json:"itemsSearchKey"`
	Labels              *[]string `json:"labels,omitempty"`
	NumberOfLinesToSkip *int64    `json:"numberOfLinesToSkip,omitempty"`
	Provider            string    `json:"provider"`
	RawContent          *string   `json:"rawContent,omitempty"`
	Source              Source    `json:"source"`
	TenantId            *string   `json:"tenantId,omitempty"`
	Updated             *string   `json:"updated,omitempty"`
	UpdatedBy           *UserInfo `json:"updatedBy,omitempty"`
	UploadStatus        *string   `json:"uploadStatus,omitempty"`
	WatchlistAlias      *string   `json:"watchlistAlias,omitempty"`
	WatchlistId         *string   `json:"watchlistId,omitempty"`
	WatchlistType       *string   `json:"watchlistType,omitempty"`
}

func (o *WatchlistProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *WatchlistProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *WatchlistProperties) GetUpdatedAsTime() (*time.Time, error) {
	if o.Updated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Updated, "2006-01-02T15:04:05Z07:00")
}

func (o *WatchlistProperties) SetUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Updated = &formatted
}
