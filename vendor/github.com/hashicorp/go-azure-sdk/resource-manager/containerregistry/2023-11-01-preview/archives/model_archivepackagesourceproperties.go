package archives

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArchivePackageSourceProperties struct {
	Type *PackageSourceType `json:"type,omitempty"`
	Url  *string            `json:"url,omitempty"`
}
