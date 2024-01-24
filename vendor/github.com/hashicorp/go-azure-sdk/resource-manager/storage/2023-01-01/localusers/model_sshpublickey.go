package localusers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SshPublicKey struct {
	Description *string `json:"description,omitempty"`
	Key         *string `json:"key,omitempty"`
}
