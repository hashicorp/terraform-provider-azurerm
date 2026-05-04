package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPSecurityRestriction struct {
	Action               *string              `json:"action,omitempty"`
	Description          *string              `json:"description,omitempty"`
	Headers              *map[string][]string `json:"headers,omitempty"`
	IPAddress            *string              `json:"ipAddress,omitempty"`
	Name                 *string              `json:"name,omitempty"`
	Priority             *int64               `json:"priority,omitempty"`
	SubnetMask           *string              `json:"subnetMask,omitempty"`
	SubnetTrafficTag     *int64               `json:"subnetTrafficTag,omitempty"`
	Tag                  *IPFilterTag         `json:"tag,omitempty"`
	VnetSubnetResourceId *string              `json:"vnetSubnetResourceId,omitempty"`
	VnetTrafficTag       *int64               `json:"vnetTrafficTag,omitempty"`
}
