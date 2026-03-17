package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobRestoreParameters struct {
	BlobRanges    []BlobRestoreRange `json:"blobRanges"`
	TimeToRestore string             `json:"timeToRestore"`
}

func (o *BlobRestoreParameters) GetTimeToRestoreAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.TimeToRestore, "2006-01-02T15:04:05Z07:00")
}

func (o *BlobRestoreParameters) SetTimeToRestoreAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeToRestore = formatted
}
