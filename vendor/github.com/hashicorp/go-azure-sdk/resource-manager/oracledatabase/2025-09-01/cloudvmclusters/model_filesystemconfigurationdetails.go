package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileSystemConfigurationDetails struct {
	FileSystemSizeGb *int64  `json:"fileSystemSizeGb,omitempty"`
	MountPoint       *string `json:"mountPoint,omitempty"`
}
