package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmSyncDetails struct {
	Last15MinutesTransferredBytes *int64                         `json:"last15MinutesTransferredBytes,omitempty"`
	LastDataTransferTimeUtc       *string                        `json:"lastDataTransferTimeUtc,omitempty"`
	LastRefreshTime               *string                        `json:"lastRefreshTime,omitempty"`
	ProcessedBytes                *int64                         `json:"processedBytes,omitempty"`
	ProgressHealth                *DiskReplicationProgressHealth `json:"progressHealth,omitempty"`
	ProgressPercentage            *int64                         `json:"progressPercentage,omitempty"`
	StartTime                     *string                        `json:"startTime,omitempty"`
	TransferredBytes              *int64                         `json:"transferredBytes,omitempty"`
}
