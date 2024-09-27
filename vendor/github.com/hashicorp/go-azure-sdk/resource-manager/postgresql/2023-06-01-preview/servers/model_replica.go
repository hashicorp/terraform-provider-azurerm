package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Replica struct {
	Capacity         *int64                    `json:"capacity,omitempty"`
	PromoteMode      *ReadReplicaPromoteMode   `json:"promoteMode,omitempty"`
	PromoteOption    *ReplicationPromoteOption `json:"promoteOption,omitempty"`
	ReplicationState *ReplicationState         `json:"replicationState,omitempty"`
	Role             *ReplicationRole          `json:"role,omitempty"`
}
