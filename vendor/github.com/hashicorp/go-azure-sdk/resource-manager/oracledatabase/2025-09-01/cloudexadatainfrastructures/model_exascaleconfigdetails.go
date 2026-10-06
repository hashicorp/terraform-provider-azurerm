package cloudexadatainfrastructures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExascaleConfigDetails struct {
	AvailableStorageInGbs *int64 `json:"availableStorageInGbs,omitempty"`
	TotalStorageInGbs     int64  `json:"totalStorageInGbs"`
}
