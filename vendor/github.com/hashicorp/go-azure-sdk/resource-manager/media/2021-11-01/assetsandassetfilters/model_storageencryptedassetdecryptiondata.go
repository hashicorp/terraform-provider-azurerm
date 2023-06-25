package assetsandassetfilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageEncryptedAssetDecryptionData struct {
	AssetFileEncryptionMetadata *[]AssetFileEncryptionMetadata `json:"assetFileEncryptionMetadata,omitempty"`
	Key                         *string                        `json:"key,omitempty"`
}
