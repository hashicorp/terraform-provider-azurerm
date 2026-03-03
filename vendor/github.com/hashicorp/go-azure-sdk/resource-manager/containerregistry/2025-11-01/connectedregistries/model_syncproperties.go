package connectedregistries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SyncProperties struct {
	GatewayEndpoint *string `json:"gatewayEndpoint,omitempty"`
	LastSyncTime    *string `json:"lastSyncTime,omitempty"`
	MessageTtl      string  `json:"messageTtl"`
	Schedule        *string `json:"schedule,omitempty"`
	SyncWindow      *string `json:"syncWindow,omitempty"`
	TokenId         string  `json:"tokenId"`
}

func (o *SyncProperties) GetLastSyncTimeAsTime() (*time.Time, error) {
	if o.LastSyncTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSyncTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SyncProperties) SetLastSyncTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSyncTime = &formatted
}
