// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_certificate_profile

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *certificateprofiles.CertificateProfileId, model *certificateprofiles.CertificateProfile) error {
	state := TrustedSigningCertificateProfileModel{
		Name:                    id.CertificateProfileName,
		TrustedSigningAccountId: codesigningaccounts.NewCodeSigningAccountID(id.SubscriptionId, id.ResourceGroupName, id.CodeSigningAccountName).ID(),
	}
	if model != nil {
		if props := model.Properties; props != nil {
			state.IdentityValidationId = props.IdentityValidationId
			state.ProfileType = string(props.ProfileType)
			state.IncludeCity = pointer.From(props.IncludeCity)
			state.IncludeCountry = pointer.From(props.IncludeCountry)
			state.IncludePostalCode = pointer.From(props.IncludePostalCode)
			state.IncludeState = pointer.From(props.IncludeState)
			state.IncludeStreetAddress = pointer.From(props.IncludeStreetAddress)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
