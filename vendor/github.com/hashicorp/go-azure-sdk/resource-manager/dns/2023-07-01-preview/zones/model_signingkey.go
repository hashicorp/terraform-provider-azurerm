package zones

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SigningKey struct {
	DelegationSignerInfo  *[]DelegationSignerInfo `json:"delegationSignerInfo,omitempty"`
	Flags                 *int64                  `json:"flags,omitempty"`
	KeyTag                *int64                  `json:"keyTag,omitempty"`
	Protocol              *int64                  `json:"protocol,omitempty"`
	PublicKey             *string                 `json:"publicKey,omitempty"`
	SecurityAlgorithmType *int64                  `json:"securityAlgorithmType,omitempty"`
}
