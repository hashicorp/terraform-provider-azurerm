package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinuxOperatingSystemProfile struct {
	Password   *string     `json:"password,omitempty"`
	SshProfile *SshProfile `json:"sshProfile,omitempty"`
	Username   *string     `json:"username,omitempty"`
}
