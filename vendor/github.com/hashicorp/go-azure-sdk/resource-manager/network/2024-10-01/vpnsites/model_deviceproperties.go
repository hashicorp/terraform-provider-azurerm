package vpnsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeviceProperties struct {
	DeviceModel     *string `json:"deviceModel,omitempty"`
	DeviceVendor    *string `json:"deviceVendor,omitempty"`
	LinkSpeedInMbps *int64  `json:"linkSpeedInMbps,omitempty"`
}
