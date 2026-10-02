package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyReleaseParameters struct {
	Enc    *KeyEncryptionAlgorithm `json:"enc,omitempty"`
	Nonce  *string                 `json:"nonce,omitempty"`
	Target string                  `json:"target"`
}
