package sapvirtualinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SshKeyPair struct {
	PrivateKey *string `json:"privateKey,omitempty"`
	PublicKey  *string `json:"publicKey,omitempty"`
}
