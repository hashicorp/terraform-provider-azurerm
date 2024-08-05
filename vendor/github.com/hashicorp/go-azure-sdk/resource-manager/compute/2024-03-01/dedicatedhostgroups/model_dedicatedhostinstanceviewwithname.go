package dedicatedhostgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostInstanceViewWithName struct {
	AssetId           *string                         `json:"assetId,omitempty"`
	AvailableCapacity *DedicatedHostAvailableCapacity `json:"availableCapacity,omitempty"`
	Name              *string                         `json:"name,omitempty"`
	Statuses          *[]InstanceViewStatus           `json:"statuses,omitempty"`
}
