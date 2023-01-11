package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Snapshot struct {
	Id         *string             `json:"id,omitempty"`
	Location   string              `json:"location"`
	Name       *string             `json:"name,omitempty"`
	Properties *SnapshotProperties `json:"properties,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
