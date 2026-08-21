package loadbalancers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointIPConfigurationProperties struct {
	GroupId          *string `json:"groupId,omitempty"`
	MemberName       *string `json:"memberName,omitempty"`
	PrivateIPAddress *string `json:"privateIPAddress,omitempty"`
}
