package openapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DelegatedSubnetUsage struct {
	SubnetName *string `json:"subnetName,omitempty"`
	Usage      *int64  `json:"usage,omitempty"`
}
