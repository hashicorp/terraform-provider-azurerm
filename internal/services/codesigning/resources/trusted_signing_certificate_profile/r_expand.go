// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_certificate_profile

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
)

func (Resource) expand(model TrustedSigningCertificateProfileModel) certificateprofiles.CertificateProfile {
	return certificateprofiles.CertificateProfile{
		Properties: &certificateprofiles.CertificateProfileProperties{
			IdentityValidationId: model.IdentityValidationId,
			ProfileType:          certificateprofiles.ProfileType(model.ProfileType),
			IncludeCity:          pointer.To(model.IncludeCity),
			IncludeCountry:       pointer.To(model.IncludeCountry),
			IncludePostalCode:    pointer.To(model.IncludePostalCode),
			IncludeState:         pointer.To(model.IncludeState),
			IncludeStreetAddress: pointer.To(model.IncludeStreetAddress),
		},
	}
}
