package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteFailoverCircuitResourceDetails struct {
	ConnectionName *string `json:"connectionName,omitempty"`
	Name           *string `json:"name,omitempty"`
	NrpResourceUri *string `json:"nrpResourceUri,omitempty"`
}
