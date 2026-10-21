package bigdatapools

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LibraryInfo struct {
	ContainerName      *string `json:"containerName,omitempty"`
	CreatorId          *string `json:"creatorId,omitempty"`
	Name               *string `json:"name,omitempty"`
	Path               *string `json:"path,omitempty"`
	ProvisioningStatus *string `json:"provisioningStatus,omitempty"`
	Type               *string `json:"type,omitempty"`
	UploadedTimestamp  *string `json:"uploadedTimestamp,omitempty"`
}

func (o *LibraryInfo) GetUploadedTimestampAsTime() (*time.Time, error) {
	if o.UploadedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UploadedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *LibraryInfo) SetUploadedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UploadedTimestamp = &formatted
}
