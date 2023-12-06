package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RequiredTraffic struct {
	Direction *TrafficDirection `json:"direction,omitempty"`
	Fqdns     *[]string         `json:"fqdns,omitempty"`
	IPs       *[]string         `json:"ips,omitempty"`
	Port      *int64            `json:"port,omitempty"`
	Protocol  *string           `json:"protocol,omitempty"`
}
