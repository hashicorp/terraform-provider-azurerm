package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnLinkProviderProperties struct {
	LinkProviderName *string `json:"linkProviderName,omitempty"`
	LinkSpeedInMbps  *int64  `json:"linkSpeedInMbps,omitempty"`
}
