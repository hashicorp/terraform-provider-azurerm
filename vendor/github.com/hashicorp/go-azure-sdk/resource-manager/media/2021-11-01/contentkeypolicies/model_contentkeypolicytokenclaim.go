package contentkeypolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPolicyTokenClaim struct {
	ClaimType  *string `json:"claimType,omitempty"`
	ClaimValue *string `json:"claimValue,omitempty"`
}
