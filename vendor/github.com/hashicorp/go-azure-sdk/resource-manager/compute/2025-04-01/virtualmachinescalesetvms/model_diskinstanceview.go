package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskInstanceView struct {
	EncryptionSettings *[]DiskEncryptionSettings `json:"encryptionSettings,omitempty"`
	Name               *string                   `json:"name,omitempty"`
	Statuses           *[]InstanceViewStatus     `json:"statuses,omitempty"`
}
