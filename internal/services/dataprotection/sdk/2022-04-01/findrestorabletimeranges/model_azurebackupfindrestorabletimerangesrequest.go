package findrestorabletimeranges

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBackupFindRestorableTimeRangesRequest struct {
	EndTime             *string                    `json:"endTime,omitempty"`
	SourceDataStoreType RestoreSourceDataStoreType `json:"sourceDataStoreType"`
	StartTime           *string                    `json:"startTime,omitempty"`
}
