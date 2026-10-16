// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package codesigning

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -properties name -compare-values subscription_id:trusted_signing_account_id,resource_group_name:trusted_signing_account_id,code_signing_account_name:trusted_signing_account_id -test-env-vars ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID

type TrustedSigningCertificateProfileResource struct{}

var _ sdk.ResourceWithIdentity = TrustedSigningCertificateProfileResource{}

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

func (TrustedSigningCertificateProfileResource) ResourceType() string {
	return "azurerm_trusted_signing_certificate_profile"
}

func (TrustedSigningCertificateProfileResource) ModelObject() interface{} {
	return &TrustedSigningCertificateProfileModel{}
}

func (TrustedSigningCertificateProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificateprofiles.ValidateCertificateProfileID
}

func (TrustedSigningCertificateProfileResource) Identity() resourceids.ResourceId {
	return &certificateprofiles.CertificateProfileId{}
}

func (TrustedSigningCertificateProfileResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(5, 100),
				validation.StringMatch(regexp.MustCompile("^[A-Za-z][A-Za-z0-9]*(-[A-Za-z0-9]+)*$"), "must begin with a letter, end with a letter or digit, and contain only alphanumeric characters and non-consecutive hyphens"),
			),
		},

		"identity_validation_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"profile_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(certificateprofiles.PossibleValuesForProfileType(), false),
		},

		"trusted_signing_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: codesigningaccounts.ValidateCodeSigningAccountID,
		},

		"include_city": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"include_country": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"include_postal_code": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"include_state": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"include_street_address": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
	}
}

func (TrustedSigningCertificateProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r TrustedSigningCertificateProfileResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CodeSigning.Client.CertificateProfiles
			var model TrustedSigningCertificateProfileModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := codesigningaccounts.ParseCodeSigningAccountID(model.TrustedSigningAccountId)
			if err != nil {
				return err
			}
			id := certificateprofiles.NewCertificateProfileID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.CodeSigningAccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			if err := client.CreateCallbackThenPoll(ctx, id, r.expand(model), metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r TrustedSigningCertificateProfileResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CodeSigning.Client.CertificateProfiles
			id, err := certificateprofiles.ParseCertificateProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (TrustedSigningCertificateProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CodeSigning.Client.CertificateProfiles
			id, err := certificateprofiles.ParseCertificateProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (TrustedSigningCertificateProfileResource) expand(model TrustedSigningCertificateProfileModel) certificateprofiles.CertificateProfile {
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

func (TrustedSigningCertificateProfileResource) flatten(metadata sdk.ResourceMetaData, id *certificateprofiles.CertificateProfileId, model *certificateprofiles.CertificateProfile) error {
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
