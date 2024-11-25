package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskContainerSettings struct {
	ContainerHostBatchBindMounts *[]ContainerHostBatchBindMountEntry `json:"containerHostBatchBindMounts,omitempty"`
	ContainerRunOptions          *string                             `json:"containerRunOptions,omitempty"`
	ImageName                    string                              `json:"imageName"`
	Registry                     *ContainerRegistry                  `json:"registry,omitempty"`
	WorkingDirectory             *ContainerWorkingDirectory          `json:"workingDirectory,omitempty"`
}
