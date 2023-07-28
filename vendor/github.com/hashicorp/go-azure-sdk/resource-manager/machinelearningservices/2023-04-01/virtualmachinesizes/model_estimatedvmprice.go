package virtualmachinesizes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EstimatedVMPrice struct {
	OsType      VMPriceOSType `json:"osType"`
	RetailPrice float64       `json:"retailPrice"`
	VMTier      VMTier        `json:"vmTier"`
}
