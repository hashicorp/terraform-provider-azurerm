package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InboundNatPool struct {
	BackendPort               int64                       `json:"backendPort"`
	FrontendPortRangeEnd      int64                       `json:"frontendPortRangeEnd"`
	FrontendPortRangeStart    int64                       `json:"frontendPortRangeStart"`
	Name                      string                      `json:"name"`
	NetworkSecurityGroupRules *[]NetworkSecurityGroupRule `json:"networkSecurityGroupRules,omitempty"`
	Protocol                  InboundEndpointProtocol     `json:"protocol"`
}
