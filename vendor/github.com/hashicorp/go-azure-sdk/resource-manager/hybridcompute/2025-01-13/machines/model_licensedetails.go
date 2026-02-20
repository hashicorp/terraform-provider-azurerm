package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseDetails struct {
	AssignedLicenses     *int64                  `json:"assignedLicenses,omitempty"`
	Edition              *LicenseEdition         `json:"edition,omitempty"`
	ImmutableId          *string                 `json:"immutableId,omitempty"`
	Processors           *int64                  `json:"processors,omitempty"`
	State                *LicenseState           `json:"state,omitempty"`
	Target               *LicenseTarget          `json:"target,omitempty"`
	Type                 *LicenseCoreType        `json:"type,omitempty"`
	VolumeLicenseDetails *[]VolumeLicenseDetails `json:"volumeLicenseDetails,omitempty"`
}
