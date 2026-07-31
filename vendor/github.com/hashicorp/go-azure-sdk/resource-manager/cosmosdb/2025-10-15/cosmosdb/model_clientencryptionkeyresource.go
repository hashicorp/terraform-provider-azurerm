package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientEncryptionKeyResource struct {
	EncryptionAlgorithm      *string          `json:"encryptionAlgorithm,omitempty"`
	Id                       *string          `json:"id,omitempty"`
	KeyWrapMetadata          *KeyWrapMetadata `json:"keyWrapMetadata,omitempty"`
	WrappedDataEncryptionKey *string          `json:"wrappedDataEncryptionKey,omitempty"`
}
