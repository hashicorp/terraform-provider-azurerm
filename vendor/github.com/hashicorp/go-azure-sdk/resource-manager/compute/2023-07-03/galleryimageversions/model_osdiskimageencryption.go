package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSDiskImageEncryption struct {
	DiskEncryptionSetId *string                     `json:"diskEncryptionSetId,omitempty"`
	SecurityProfile     *OSDiskImageSecurityProfile `json:"securityProfile,omitempty"`
}
