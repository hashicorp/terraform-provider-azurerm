package brokerlistener

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertManagerPrivateKey struct {
	Algorithm      PrivateKeyAlgorithm      `json:"algorithm"`
	RotationPolicy PrivateKeyRotationPolicy `json:"rotationPolicy"`
}
