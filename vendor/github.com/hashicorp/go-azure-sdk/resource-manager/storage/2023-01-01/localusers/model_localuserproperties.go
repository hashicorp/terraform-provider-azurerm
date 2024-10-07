package localusers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalUserProperties struct {
	HasSharedKey      *bool              `json:"hasSharedKey,omitempty"`
	HasSshKey         *bool              `json:"hasSshKey,omitempty"`
	HasSshPassword    *bool              `json:"hasSshPassword,omitempty"`
	HomeDirectory     *string            `json:"homeDirectory,omitempty"`
	PermissionScopes  *[]PermissionScope `json:"permissionScopes,omitempty"`
	Sid               *string            `json:"sid,omitempty"`
	SshAuthorizedKeys *[]SshPublicKey    `json:"sshAuthorizedKeys,omitempty"`
}
