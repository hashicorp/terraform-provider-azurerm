package orders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrackingInfo struct {
	CarrierName  *string `json:"carrierName,omitempty"`
	SerialNumber *string `json:"serialNumber,omitempty"`
	TrackingId   *string `json:"trackingId,omitempty"`
	TrackingUrl  *string `json:"trackingUrl,omitempty"`
}
