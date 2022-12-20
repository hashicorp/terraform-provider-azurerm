package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotRestoreFiles struct {
	DestinationPath *string  `json:"destinationPath,omitempty"`
	FilePaths       []string `json:"filePaths"`
}
