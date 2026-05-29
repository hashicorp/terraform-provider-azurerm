package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientEncryptionIncludedPath struct {
	ClientEncryptionKeyId string `json:"clientEncryptionKeyId"`
	EncryptionAlgorithm   string `json:"encryptionAlgorithm"`
	EncryptionType        string `json:"encryptionType"`
	Path                  string `json:"path"`
}
