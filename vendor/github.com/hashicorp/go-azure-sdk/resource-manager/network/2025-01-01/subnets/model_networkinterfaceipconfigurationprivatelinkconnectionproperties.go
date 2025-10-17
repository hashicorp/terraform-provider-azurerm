package subnets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfaceIPConfigurationPrivateLinkConnectionProperties struct {
	Fqdns              *[]string `json:"fqdns,omitempty"`
	GroupId            *string   `json:"groupId,omitempty"`
	RequiredMemberName *string   `json:"requiredMemberName,omitempty"`
}
