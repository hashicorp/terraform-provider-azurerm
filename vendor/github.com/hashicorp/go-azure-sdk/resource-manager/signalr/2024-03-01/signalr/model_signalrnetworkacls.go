package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SignalRNetworkACLs struct {
	DefaultAction    *ACLAction            `json:"defaultAction,omitempty"`
	IPRules          *[]IPRule             `json:"ipRules,omitempty"`
	PrivateEndpoints *[]PrivateEndpointACL `json:"privateEndpoints,omitempty"`
	PublicNetwork    *NetworkACL           `json:"publicNetwork,omitempty"`
}
