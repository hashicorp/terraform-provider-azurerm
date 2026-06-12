package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprAppHealth struct {
	Enabled                  *bool   `json:"enabled,omitempty"`
	Path                     *string `json:"path,omitempty"`
	ProbeIntervalSeconds     *int64  `json:"probeIntervalSeconds,omitempty"`
	ProbeTimeoutMilliseconds *int64  `json:"probeTimeoutMilliseconds,omitempty"`
	Threshold                *int64  `json:"threshold,omitempty"`
}
