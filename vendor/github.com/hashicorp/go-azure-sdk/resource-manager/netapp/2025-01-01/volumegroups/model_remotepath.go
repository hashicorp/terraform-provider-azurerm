package volumegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemotePath struct {
	ExternalHostName string `json:"externalHostName"`
	ServerName       string `json:"serverName"`
	VolumeName       string `json:"volumeName"`
}
