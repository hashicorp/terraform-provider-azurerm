package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UefiKeySignatures struct {
	Db  *[]UefiKey `json:"db,omitempty"`
	Dbx *[]UefiKey `json:"dbx,omitempty"`
	Kek *[]UefiKey `json:"kek,omitempty"`
	Pk  *UefiKey   `json:"pk,omitempty"`
}
