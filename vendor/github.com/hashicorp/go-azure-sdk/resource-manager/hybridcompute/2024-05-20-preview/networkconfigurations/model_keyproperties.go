package networkconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyProperties struct {
	CandidatePublicKey *KeyDetails `json:"candidatePublicKey,omitempty"`
	ClientPublicKey    *KeyDetails `json:"clientPublicKey,omitempty"`
}
