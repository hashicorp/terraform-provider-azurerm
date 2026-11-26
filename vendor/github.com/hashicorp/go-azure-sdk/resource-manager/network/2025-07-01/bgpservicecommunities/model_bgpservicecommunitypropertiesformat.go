package bgpservicecommunities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BgpServiceCommunityPropertiesFormat struct {
	BgpCommunities *[]BGPCommunity `json:"bgpCommunities,omitempty"`
	ServiceName    *string         `json:"serviceName,omitempty"`
}
