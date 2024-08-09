package licenseprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EsuKey struct {
	LicenseStatus *int64  `json:"licenseStatus,omitempty"`
	Sku           *string `json:"sku,omitempty"`
}
