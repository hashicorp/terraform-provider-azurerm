package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Checkpoint struct {
	CheckpointID       *string `json:"checkpointID,omitempty"`
	Description        *string `json:"description,omitempty"`
	Name               *string `json:"name,omitempty"`
	ParentCheckpointID *string `json:"parentCheckpointID,omitempty"`
}
