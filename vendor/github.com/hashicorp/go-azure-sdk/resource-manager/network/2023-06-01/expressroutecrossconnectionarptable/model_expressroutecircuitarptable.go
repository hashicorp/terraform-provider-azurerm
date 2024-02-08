package expressroutecrossconnectionarptable

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitArpTable struct {
	Age        *int64  `json:"age,omitempty"`
	IPAddress  *string `json:"ipAddress,omitempty"`
	Interface  *string `json:"interface,omitempty"`
	MacAddress *string `json:"macAddress,omitempty"`
}
