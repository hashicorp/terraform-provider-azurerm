package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSDiskImageSecurityProfile struct {
	ConfidentialVMEncryptionType *ConfidentialVMEncryptionType `json:"confidentialVMEncryptionType,omitempty"`
	SecureVMDiskEncryptionSetId  *string                       `json:"secureVMDiskEncryptionSetId,omitempty"`
}
