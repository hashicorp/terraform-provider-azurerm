package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MountTargetProperties struct {
	FileSystemId  string  `json:"fileSystemId"`
	IPAddress     *string `json:"ipAddress,omitempty"`
	MountTargetId *string `json:"mountTargetId,omitempty"`
	SmbServerFqdn *string `json:"smbServerFqdn,omitempty"`
}
