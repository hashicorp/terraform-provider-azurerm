package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayURLConfiguration struct {
	ModifiedPath        *string `json:"modifiedPath,omitempty"`
	ModifiedQueryString *string `json:"modifiedQueryString,omitempty"`
	Reroute             *bool   `json:"reroute,omitempty"`
}
