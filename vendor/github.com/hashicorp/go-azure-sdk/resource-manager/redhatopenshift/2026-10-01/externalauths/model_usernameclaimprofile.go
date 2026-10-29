package externalauths

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsernameClaimProfile struct {
	Claim        string                     `json:"claim"`
	Prefix       *string                    `json:"prefix,omitempty"`
	PrefixPolicy *UsernameClaimPrefixPolicy `json:"prefixPolicy,omitempty"`
}
