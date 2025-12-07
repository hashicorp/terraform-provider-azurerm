package mongoclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PromoteReplicaRequest struct {
	Mode          *PromoteMode  `json:"mode,omitempty"`
	PromoteOption PromoteOption `json:"promoteOption"`
}
