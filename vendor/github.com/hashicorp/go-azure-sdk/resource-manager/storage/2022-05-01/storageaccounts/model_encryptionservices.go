package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionServices struct {
	Blob  *EncryptionService `json:"blob"`
	File  *EncryptionService `json:"file"`
	Queue *EncryptionService `json:"queue"`
	Table *EncryptionService `json:"table"`
}
