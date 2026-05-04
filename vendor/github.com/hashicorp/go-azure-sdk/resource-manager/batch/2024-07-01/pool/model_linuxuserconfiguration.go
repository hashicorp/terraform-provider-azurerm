package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinuxUserConfiguration struct {
	Gid           *int64  `json:"gid,omitempty"`
	SshPrivateKey *string `json:"sshPrivateKey,omitempty"`
	Uid           *int64  `json:"uid,omitempty"`
}
