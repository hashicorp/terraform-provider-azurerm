package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnConnectionPacketCaptureStartParameters struct {
	FilterData          *string   `json:"filterData,omitempty"`
	LinkConnectionNames *[]string `json:"linkConnectionNames,omitempty"`
}
