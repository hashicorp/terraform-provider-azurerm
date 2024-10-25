package diskencryptionsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyForDiskEncryptionSet struct {
	KeyURL      string       `json:"keyUrl"`
	SourceVault *SourceVault `json:"sourceVault,omitempty"`
}
