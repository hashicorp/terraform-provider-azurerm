package labs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageReference struct {
	Offer     *string `json:"offer,omitempty"`
	OsType    *string `json:"osType,omitempty"`
	Publisher *string `json:"publisher,omitempty"`
	Sku       *string `json:"sku,omitempty"`
	Version   *string `json:"version,omitempty"`
}
