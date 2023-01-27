package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Encryption struct {
	DiskEncryptionSetId *string         `json:"diskEncryptionSetId,omitempty"`
	Type                *EncryptionType `json:"type,omitempty"`
}
