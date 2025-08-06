package serviceresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSyncCompleteCommandInput struct {
	CommitTimeStamp *string `json:"commitTimeStamp,omitempty"`
	DatabaseName    string  `json:"databaseName"`
}

func (o *MigrateSyncCompleteCommandInput) GetCommitTimeStampAsTime() (*time.Time, error) {
	if o.CommitTimeStamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CommitTimeStamp, "2006-01-02T15:04:05Z07:00")
}

func (o *MigrateSyncCompleteCommandInput) SetCommitTimeStampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CommitTimeStamp = &formatted
}
