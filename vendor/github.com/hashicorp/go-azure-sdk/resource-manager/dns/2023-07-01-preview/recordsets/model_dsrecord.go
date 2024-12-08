package recordsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DsRecord struct {
	Algorithm *int64  `json:"algorithm,omitempty"`
	Digest    *Digest `json:"digest,omitempty"`
	KeyTag    *int64  `json:"keyTag,omitempty"`
}
