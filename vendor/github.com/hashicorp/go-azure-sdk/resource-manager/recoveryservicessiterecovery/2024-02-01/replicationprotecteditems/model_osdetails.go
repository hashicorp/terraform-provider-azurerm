package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSDetails struct {
	OSMajorVersion     *string `json:"oSMajorVersion,omitempty"`
	OSMinorVersion     *string `json:"oSMinorVersion,omitempty"`
	OSVersion          *string `json:"oSVersion,omitempty"`
	OsEdition          *string `json:"osEdition,omitempty"`
	OsType             *string `json:"osType,omitempty"`
	ProductType        *string `json:"productType,omitempty"`
	UserSelectedOSName *string `json:"userSelectedOSName,omitempty"`
}
