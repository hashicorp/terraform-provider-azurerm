package appliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSHKey struct {
	Certificate         *string `json:"certificate,omitempty"`
	CreationTimeStamp   *int64  `json:"creationTimeStamp,omitempty"`
	ExpirationTimeStamp *int64  `json:"expirationTimeStamp,omitempty"`
	PrivateKey          *string `json:"privateKey,omitempty"`
	PublicKey           *string `json:"publicKey,omitempty"`
}
