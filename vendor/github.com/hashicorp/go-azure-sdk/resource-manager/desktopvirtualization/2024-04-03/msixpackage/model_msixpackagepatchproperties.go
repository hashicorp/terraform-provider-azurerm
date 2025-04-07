package msixpackage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MSIXPackagePatchProperties struct {
	DisplayName           *string `json:"displayName,omitempty"`
	IsActive              *bool   `json:"isActive,omitempty"`
	IsRegularRegistration *bool   `json:"isRegularRegistration,omitempty"`
}
