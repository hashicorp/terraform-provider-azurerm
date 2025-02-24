package replicationfabrics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverProcessServerRequestProperties struct {
	ContainerName         *string   `json:"containerName,omitempty"`
	SourceProcessServerId *string   `json:"sourceProcessServerId,omitempty"`
	TargetProcessServerId *string   `json:"targetProcessServerId,omitempty"`
	UpdateType            *string   `json:"updateType,omitempty"`
	VMsToMigrate          *[]string `json:"vmsToMigrate,omitempty"`
}
