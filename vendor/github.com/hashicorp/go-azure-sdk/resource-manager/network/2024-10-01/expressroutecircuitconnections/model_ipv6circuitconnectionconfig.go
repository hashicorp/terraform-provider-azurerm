package expressroutecircuitconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPv6CircuitConnectionConfig struct {
	AddressPrefix           *string                  `json:"addressPrefix,omitempty"`
	CircuitConnectionStatus *CircuitConnectionStatus `json:"circuitConnectionStatus,omitempty"`
}
