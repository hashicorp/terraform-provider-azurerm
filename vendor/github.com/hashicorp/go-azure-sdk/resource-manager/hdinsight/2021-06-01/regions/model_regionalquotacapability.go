package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegionalQuotaCapability struct {
	CoresAvailable *int64  `json:"coresAvailable,omitempty"`
	CoresUsed      *int64  `json:"coresUsed,omitempty"`
	RegionName     *string `json:"regionName,omitempty"`
}
