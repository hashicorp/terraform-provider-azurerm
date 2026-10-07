package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionServices struct {
	Blob  *EncryptionService `json:"blob,omitempty"`
	File  *EncryptionService `json:"file,omitempty"`
	Queue *EncryptionService `json:"queue,omitempty"`
	Table *EncryptionService `json:"table,omitempty"`
}
