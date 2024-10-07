package resolveprivatelinkserviceid

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkResource struct {
	GroupId              *string   `json:"groupId,omitempty"`
	Id                   *string   `json:"id,omitempty"`
	Name                 *string   `json:"name,omitempty"`
	PrivateLinkServiceID *string   `json:"privateLinkServiceID,omitempty"`
	RequiredMembers      *[]string `json:"requiredMembers,omitempty"`
	Type                 *string   `json:"type,omitempty"`
}
