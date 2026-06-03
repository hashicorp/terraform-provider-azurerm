package fileservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NfsSetting struct {
	EncryptionInTransit *EncryptionInTransit `json:"encryptionInTransit,omitempty"`
}
