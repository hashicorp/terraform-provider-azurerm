package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SshPublicKey struct {
	KeyData *string `json:"keyData,omitempty"`
	Path    *string `json:"path,omitempty"`
}
