package galleryapplicationversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationCustomActionParameter struct {
	DefaultValue *string                                      `json:"defaultValue,omitempty"`
	Description  *string                                      `json:"description,omitempty"`
	Name         string                                       `json:"name"`
	Required     *bool                                        `json:"required,omitempty"`
	Type         *GalleryApplicationCustomActionParameterType `json:"type,omitempty"`
}
