// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_certificate_profile

type TrustedSigningCertificateProfileModel struct {
	Name                    string `tfschema:"name"`
	IdentityValidationId    string `tfschema:"identity_validation_id"`
	ProfileType             string `tfschema:"profile_type"`
	TrustedSigningAccountId string `tfschema:"trusted_signing_account_id"`
	IncludeCity             bool   `tfschema:"include_city"`
	IncludeCountry          bool   `tfschema:"include_country"`
	IncludePostalCode       bool   `tfschema:"include_postal_code"`
	IncludeState            bool   `tfschema:"include_state"`
	IncludeStreetAddress    bool   `tfschema:"include_street_address"`
}

func (Resource) ModelObject() any {
	return &TrustedSigningCertificateProfileModel{}
}
