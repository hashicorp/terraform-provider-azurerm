package appliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridConnectionConfig struct {
	ExpirationTime       *int64  `json:"expirationTime,omitempty"`
	HybridConnectionName *string `json:"hybridConnectionName,omitempty"`
	Relay                *string `json:"relay,omitempty"`
	Token                *string `json:"token,omitempty"`
}
