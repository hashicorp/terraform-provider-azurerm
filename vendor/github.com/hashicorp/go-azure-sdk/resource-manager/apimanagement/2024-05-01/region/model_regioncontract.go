package region

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegionContract struct {
	IsDeleted      *bool   `json:"isDeleted,omitempty"`
	IsMasterRegion *bool   `json:"isMasterRegion,omitempty"`
	Name           *string `json:"name,omitempty"`
}
