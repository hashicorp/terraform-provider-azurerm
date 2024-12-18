package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicaExecutionStatus struct {
	Containers *[]ContainerExecutionStatus `json:"containers,omitempty"`
	Name       *string                     `json:"name,omitempty"`
}
