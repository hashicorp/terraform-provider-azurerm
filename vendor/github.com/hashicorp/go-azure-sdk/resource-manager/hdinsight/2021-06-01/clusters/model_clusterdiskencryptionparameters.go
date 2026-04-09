package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterDiskEncryptionParameters struct {
	KeyName    *string `json:"keyName,omitempty"`
	KeyVersion *string `json:"keyVersion,omitempty"`
	VaultUri   *string `json:"vaultUri,omitempty"`
}
