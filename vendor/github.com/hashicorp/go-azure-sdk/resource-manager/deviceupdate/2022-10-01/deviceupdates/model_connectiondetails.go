package deviceupdates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionDetails struct {
	GroupId          *string `json:"groupId,omitempty"`
	Id               *string `json:"id,omitempty"`
	LinkIdentifier   *string `json:"linkIdentifier,omitempty"`
	MemberName       *string `json:"memberName,omitempty"`
	PrivateIPAddress *string `json:"privateIpAddress,omitempty"`
}
