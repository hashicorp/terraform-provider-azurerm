package dbversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbVersionProperties struct {
	IsLatestForMajorVersion *bool  `json:"isLatestForMajorVersion,omitempty"`
	IsPreviewDbVersion      *bool  `json:"isPreviewDbVersion,omitempty"`
	IsUpgradeSupported      *bool  `json:"isUpgradeSupported,omitempty"`
	SupportsPdb             *bool  `json:"supportsPdb,omitempty"`
	Version                 string `json:"version"`
}
