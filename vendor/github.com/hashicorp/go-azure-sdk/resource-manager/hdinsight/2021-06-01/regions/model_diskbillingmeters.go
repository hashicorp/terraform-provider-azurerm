package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskBillingMeters struct {
	DiskRpMeter *string `json:"diskRpMeter,omitempty"`
	Sku         *string `json:"sku,omitempty"`
	Tier        *Tier   `json:"tier,omitempty"`
}
