package dedicatedhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostInstanceView struct {
	AssetId           *string                         `json:"assetId,omitempty"`
	AvailableCapacity *DedicatedHostAvailableCapacity `json:"availableCapacity,omitempty"`
	Statuses          *[]InstanceViewStatus           `json:"statuses,omitempty"`
}
