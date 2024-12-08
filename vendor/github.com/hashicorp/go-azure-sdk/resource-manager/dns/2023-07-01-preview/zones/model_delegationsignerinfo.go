package zones

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DelegationSignerInfo struct {
	DigestAlgorithmType *int64  `json:"digestAlgorithmType,omitempty"`
	DigestValue         *string `json:"digestValue,omitempty"`
	Record              *string `json:"record,omitempty"`
}
