package amlfilesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemClientInfo struct {
	ContainerStorageInterface *AmlFilesystemContainerStorageInterface `json:"containerStorageInterface,omitempty"`
	LustreVersion             *string                                 `json:"lustreVersion,omitempty"`
	MgsAddress                *string                                 `json:"mgsAddress,omitempty"`
	MountCommand              *string                                 `json:"mountCommand,omitempty"`
}
