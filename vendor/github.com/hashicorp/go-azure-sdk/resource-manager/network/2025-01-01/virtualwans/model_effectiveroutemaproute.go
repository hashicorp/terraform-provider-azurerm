package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EffectiveRouteMapRoute struct {
	AsPath         *string `json:"asPath,omitempty"`
	BgpCommunities *string `json:"bgpCommunities,omitempty"`
	Prefix         *string `json:"prefix,omitempty"`
}
