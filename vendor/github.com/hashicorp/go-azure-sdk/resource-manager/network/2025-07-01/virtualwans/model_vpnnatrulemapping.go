package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnNatRuleMapping struct {
	AddressSpace *string `json:"addressSpace,omitempty"`
	PortRange    *string `json:"portRange,omitempty"`
}
