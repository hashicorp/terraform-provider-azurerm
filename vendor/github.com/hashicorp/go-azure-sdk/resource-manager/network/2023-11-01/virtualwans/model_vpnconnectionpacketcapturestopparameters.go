package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnConnectionPacketCaptureStopParameters struct {
	LinkConnectionNames *[]string `json:"linkConnectionNames,omitempty"`
	SasUrl              *string   `json:"sasUrl,omitempty"`
}
