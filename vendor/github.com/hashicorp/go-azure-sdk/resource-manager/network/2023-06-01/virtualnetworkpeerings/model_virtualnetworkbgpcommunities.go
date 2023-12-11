package virtualnetworkpeerings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkBgpCommunities struct {
	RegionalCommunity       *string `json:"regionalCommunity,omitempty"`
	VirtualNetworkCommunity string  `json:"virtualNetworkCommunity"`
}
