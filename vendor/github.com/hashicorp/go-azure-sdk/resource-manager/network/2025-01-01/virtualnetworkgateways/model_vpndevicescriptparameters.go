package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnDeviceScriptParameters struct {
	DeviceFamily    *string `json:"deviceFamily,omitempty"`
	FirmwareVersion *string `json:"firmwareVersion,omitempty"`
	Vendor          *string `json:"vendor,omitempty"`
}
