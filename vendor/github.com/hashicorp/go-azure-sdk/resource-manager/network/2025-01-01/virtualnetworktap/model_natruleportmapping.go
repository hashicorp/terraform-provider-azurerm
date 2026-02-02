package virtualnetworktap

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NatRulePortMapping struct {
	BackendPort        *int64  `json:"backendPort,omitempty"`
	FrontendPort       *int64  `json:"frontendPort,omitempty"`
	InboundNatRuleName *string `json:"inboundNatRuleName,omitempty"`
}
