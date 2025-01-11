package virtualnetworkgatewayconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficSelectorPolicy struct {
	LocalAddressRanges  []string `json:"localAddressRanges"`
	RemoteAddressRanges []string `json:"remoteAddressRanges"`
}
