package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UefiKey struct {
	Type  *UefiKeyType `json:"type,omitempty"`
	Value *[]string    `json:"value,omitempty"`
}
