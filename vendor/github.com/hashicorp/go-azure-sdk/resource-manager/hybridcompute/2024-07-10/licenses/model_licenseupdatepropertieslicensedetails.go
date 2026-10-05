package licenses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseUpdatePropertiesLicenseDetails struct {
	Edition    *LicenseEdition  `json:"edition,omitempty"`
	Processors *int64           `json:"processors,omitempty"`
	State      *LicenseState    `json:"state,omitempty"`
	Target     *LicenseTarget   `json:"target,omitempty"`
	Type       *LicenseCoreType `json:"type,omitempty"`
}
