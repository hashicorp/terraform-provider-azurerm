package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VersionSpec struct {
	ComponentVersions *map[string]string `json:"componentVersions,omitempty"`
	DisplayName       *string            `json:"displayName,omitempty"`
	FriendlyName      *string            `json:"friendlyName,omitempty"`
	IsDefault         *bool              `json:"isDefault,omitempty"`
}
