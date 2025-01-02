package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InitialReplicationDetails struct {
	InitialReplicationProgressPercentage *string `json:"initialReplicationProgressPercentage,omitempty"`
	InitialReplicationType               *string `json:"initialReplicationType,omitempty"`
}
