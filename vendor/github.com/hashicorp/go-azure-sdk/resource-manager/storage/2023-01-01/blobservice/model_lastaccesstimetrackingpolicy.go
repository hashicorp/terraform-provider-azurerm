package blobservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LastAccessTimeTrackingPolicy struct {
	BlobType                  *[]string `json:"blobType,omitempty"`
	Enable                    bool      `json:"enable"`
	Name                      *Name     `json:"name,omitempty"`
	TrackingGranularityInDays *int64    `json:"trackingGranularityInDays,omitempty"`
}
