package attacheddatanetwork

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NaptConfiguration struct {
	Enabled           *NaptEnabled        `json:"enabled,omitempty"`
	PinholeLimits     *int64              `json:"pinholeLimits,omitempty"`
	PinholeTimeouts   *PinholeTimeouts    `json:"pinholeTimeouts,omitempty"`
	PortRange         *PortRange          `json:"portRange,omitempty"`
	PortReuseHoldTime *PortReuseHoldTimes `json:"portReuseHoldTime,omitempty"`
}
