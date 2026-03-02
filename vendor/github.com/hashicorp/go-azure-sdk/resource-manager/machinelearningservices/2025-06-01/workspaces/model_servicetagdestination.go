package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceTagDestination struct {
	Action          *RuleAction `json:"action,omitempty"`
	AddressPrefixes *[]string   `json:"addressPrefixes,omitempty"`
	PortRanges      *string     `json:"portRanges,omitempty"`
	Protocol        *string     `json:"protocol,omitempty"`
	ServiceTag      *string     `json:"serviceTag,omitempty"`
}
