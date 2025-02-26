package msixpackage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MSIXPackagePatch struct {
	Id         *string                     `json:"id,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	Properties *MSIXPackagePatchProperties `json:"properties,omitempty"`
	Type       *string                     `json:"type,omitempty"`
}
