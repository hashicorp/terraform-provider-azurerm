package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceSshSettings struct {
	AdminPublicKey  *string          `json:"adminPublicKey,omitempty"`
	AdminUserName   *string          `json:"adminUserName,omitempty"`
	SshPort         *int64           `json:"sshPort,omitempty"`
	SshPublicAccess *SshPublicAccess `json:"sshPublicAccess,omitempty"`
}
