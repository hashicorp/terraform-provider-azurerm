package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageProfile struct {
	NfsCsiDriver *StorageProfileNfsCSIDriver `json:"nfsCsiDriver,omitempty"`
	SmbCsiDriver *StorageProfileSmbCSIDriver `json:"smbCsiDriver,omitempty"`
}
