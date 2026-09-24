package bgpservicecommunities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BGPCommunity struct {
	CommunityName          *string   `json:"communityName,omitempty"`
	CommunityPrefixes      *[]string `json:"communityPrefixes,omitempty"`
	CommunityValue         *string   `json:"communityValue,omitempty"`
	IsAuthorizedToUse      *bool     `json:"isAuthorizedToUse,omitempty"`
	ServiceGroup           *string   `json:"serviceGroup,omitempty"`
	ServiceSupportedRegion *string   `json:"serviceSupportedRegion,omitempty"`
}
