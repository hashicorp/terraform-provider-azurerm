package recordsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TlsaRecord struct {
	CertAssociationData *string `json:"certAssociationData,omitempty"`
	MatchingType        *int64  `json:"matchingType,omitempty"`
	Selector            *int64  `json:"selector,omitempty"`
	Usage               *int64  `json:"usage,omitempty"`
}
