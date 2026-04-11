package certificateprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProfileProperties struct {
	Certificates         *[]Certificate            `json:"certificates,omitempty"`
	IdentityValidationId string                    `json:"identityValidationId"`
	IncludeCity          *bool                     `json:"includeCity,omitempty"`
	IncludeCountry       *bool                     `json:"includeCountry,omitempty"`
	IncludePostalCode    *bool                     `json:"includePostalCode,omitempty"`
	IncludeState         *bool                     `json:"includeState,omitempty"`
	IncludeStreetAddress *bool                     `json:"includeStreetAddress,omitempty"`
	ProfileType          ProfileType               `json:"profileType"`
	ProvisioningState    *ProvisioningState        `json:"provisioningState,omitempty"`
	Status               *CertificateProfileStatus `json:"status,omitempty"`
}
