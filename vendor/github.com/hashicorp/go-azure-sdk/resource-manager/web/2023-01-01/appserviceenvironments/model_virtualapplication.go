package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualApplication struct {
	PhysicalPath       *string             `json:"physicalPath,omitempty"`
	PreloadEnabled     *bool               `json:"preloadEnabled,omitempty"`
	VirtualDirectories *[]VirtualDirectory `json:"virtualDirectories,omitempty"`
	VirtualPath        *string             `json:"virtualPath,omitempty"`
}
