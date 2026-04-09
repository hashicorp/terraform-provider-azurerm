package workflows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlowEndpoints struct {
	AccessEndpointIPAddresses *[]IPAddress `json:"accessEndpointIpAddresses,omitempty"`
	OutgoingIPAddresses       *[]IPAddress `json:"outgoingIpAddresses,omitempty"`
}
