package galleryapplicationversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationCustomAction struct {
	Description *string                                    `json:"description,omitempty"`
	Name        string                                     `json:"name"`
	Parameters  *[]GalleryApplicationCustomActionParameter `json:"parameters,omitempty"`
	Script      string                                     `json:"script"`
}
