package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenerateCertResponse struct {
	ExpiryTimeInUTC *string `json:"expiryTimeInUTC,omitempty"`
	PrivateKey      *string `json:"privateKey,omitempty"`
	PublicKey       *string `json:"publicKey,omitempty"`
}
