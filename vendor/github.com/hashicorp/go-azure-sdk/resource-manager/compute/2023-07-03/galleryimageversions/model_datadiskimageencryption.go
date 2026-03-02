package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataDiskImageEncryption struct {
	DiskEncryptionSetId *string `json:"diskEncryptionSetId,omitempty"`
	Lun                 int64   `json:"lun"`
}
