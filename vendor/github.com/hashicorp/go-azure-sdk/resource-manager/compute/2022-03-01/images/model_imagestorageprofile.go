package images

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageStorageProfile struct {
	DataDisks     *[]ImageDataDisk `json:"dataDisks,omitempty"`
	OsDisk        *ImageOSDisk     `json:"osDisk,omitempty"`
	ZoneResilient *bool            `json:"zoneResilient,omitempty"`
}
