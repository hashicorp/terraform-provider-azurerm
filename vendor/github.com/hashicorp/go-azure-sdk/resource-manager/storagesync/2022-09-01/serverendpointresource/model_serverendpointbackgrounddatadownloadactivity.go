package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointBackgroundDataDownloadActivity struct {
	DownloadedBytes  *int64  `json:"downloadedBytes,omitempty"`
	PercentProgress  *int64  `json:"percentProgress,omitempty"`
	StartedTimestamp *string `json:"startedTimestamp,omitempty"`
	Timestamp        *string `json:"timestamp,omitempty"`
}

func (o *ServerEndpointBackgroundDataDownloadActivity) GetStartedTimestampAsTime() (*time.Time, error) {
	if o.StartedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointBackgroundDataDownloadActivity) SetStartedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedTimestamp = &formatted
}

func (o *ServerEndpointBackgroundDataDownloadActivity) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointBackgroundDataDownloadActivity) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
