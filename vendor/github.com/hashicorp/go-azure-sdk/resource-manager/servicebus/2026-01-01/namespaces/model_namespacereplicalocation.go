package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceReplicaLocation struct {
	LocationName *string        `json:"locationName,omitempty"`
	RoleType     *GeoDRRoleType `json:"roleType,omitempty"`
}
