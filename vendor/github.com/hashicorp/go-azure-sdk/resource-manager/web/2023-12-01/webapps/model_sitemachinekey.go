package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteMachineKey struct {
	Decryption    *string `json:"decryption,omitempty"`
	DecryptionKey *string `json:"decryptionKey,omitempty"`
	Validation    *string `json:"validation,omitempty"`
	ValidationKey *string `json:"validationKey,omitempty"`
}
