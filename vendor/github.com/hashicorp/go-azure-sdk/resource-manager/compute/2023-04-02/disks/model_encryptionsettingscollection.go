package disks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionSettingsCollection struct {
	Enabled                   bool                         `json:"enabled"`
	EncryptionSettings        *[]EncryptionSettingsElement `json:"encryptionSettings,omitempty"`
	EncryptionSettingsVersion *string                      `json:"encryptionSettingsVersion,omitempty"`
}
