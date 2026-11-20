package ehnamespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceReplicaLocation struct {
	ClusterArmId *string        `json:"clusterArmId,omitempty"`
	LocationName *string        `json:"locationName,omitempty"`
	ReplicaState *string        `json:"replicaState,omitempty"`
	RoleType     *GeoDRRoleType `json:"roleType,omitempty"`
}
