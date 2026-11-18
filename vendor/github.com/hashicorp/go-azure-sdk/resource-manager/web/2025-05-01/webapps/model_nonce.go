package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Nonce struct {
	NonceExpirationInterval *string `json:"nonceExpirationInterval,omitempty"`
	ValidateNonce           *bool   `json:"validateNonce,omitempty"`
}
