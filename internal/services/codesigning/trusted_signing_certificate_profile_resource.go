package codesigning

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2024-09-30-preview/certificateprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2024-09-30-preview/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type TrustedSigningCertificateProfileModel struct {
	Name                    string                                       `tfschema:"name"`
	TrustedSigningAccountId string                                       `tfschema:"trusted_signing_account_id"`
	IdentityValidationId    string                                       `tfschema:"identity_validation_id"`
	IncludeCity             bool                                         `tfschema:"include_city"`
	IncludeCountry          bool                                         `tfschema:"include_country"`
	IncludePostalCode       bool                                         `tfschema:"include_postal_code"`
	IncludeState            bool                                         `tfschema:"include_state"`
	IncludeStreetAddress    bool                                         `tfschema:"include_street_address"`
	ProfileType             certificateprofiles.ProfileType              `tfschema:"profile_type"`
	Certificates            []CertificateModel                           `tfschema:"certificates"`
	Status                  certificateprofiles.CertificateProfileStatus `tfschema:"status"`
}

type CertificateModel struct {
	CreatedDate      string                                `tfschema:"created_date"`
	EnhancedKeyUsage string                                `tfschema:"enhanced_key_usage"`
	ExpiryDate       string                                `tfschema:"expiry_date"`
	Revocation       []RevocationModel                     `tfschema:"revocation"`
	SerialNumber     string                                `tfschema:"serial_number"`
	Status           certificateprofiles.CertificateStatus `tfschema:"status"`
	SubjectName      string                                `tfschema:"subject_name"`
	Thumbprint       string                                `tfschema:"thumbprint"`
}

type RevocationModel struct {
	EffectiveAt   string                               `tfschema:"effective_at"`
	FailureReason string                               `tfschema:"failure_reason"`
	Reason        string                               `tfschema:"reason"`
	Remarks       string                               `tfschema:"remarks"`
	RequestedAt   string                               `tfschema:"requested_at"`
	Status        certificateprofiles.RevocationStatus `tfschema:"status"`
}

type TrustedSigningCertificateProfileResource struct{}

var _ sdk.ResourceWithUpdate = TrustedSigningCertificateProfileResource{}

func (r TrustedSigningCertificateProfileResource) ResourceType() string {
	return "azurerm_trusted_signing_certificate_profile"
}

func (r TrustedSigningCertificateProfileResource) ModelObject() interface{} {
	return &TrustedSigningCertificateProfileModel{}
}

func (r TrustedSigningCertificateProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificateprofiles.ValidateCertificateProfileID
}

func (r TrustedSigningCertificateProfileResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(5, 100),
				validation.StringMatch(
					regexp.MustCompile("^[A-Za-z][A-Za-z0-9]*(?:-[A-Za-z0-9]+)*$"),
					"A certificate profile's name must be between 3-24 alphanumeric characters. The name must begin with a letter, end with a letter or digit, and not contain consecutive hyphens.",
				)),
		},

		"trusted_signing_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: codesigningaccounts.ValidateCodeSigningAccountID,
		},

		"identity_validation_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"profile_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(certificateprofiles.ProfileTypePublicTrust),
				string(certificateprofiles.ProfileTypePrivateTrust),
				string(certificateprofiles.ProfileTypePrivateTrustCIPolicy),
				string(certificateprofiles.ProfileTypeVBSEnclave),
				string(certificateprofiles.ProfileTypePublicTrustTest),
			}, false),
		},

		"include_city": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"include_country": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"include_postal_code": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"include_state": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"include_street_address": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r TrustedSigningCertificateProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"certificates": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"created_date": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"enhanced_key_usage": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"expiry_date": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"revocation": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"effective_at": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"failure_reason": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"reason": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"remarks": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"requested_at": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"status": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"serial_number": {
						Type:      pluginsdk.TypeString,
						Sensitive: true,
						Computed:  true,
					},

					"status": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"subject_name": {
						Type:      pluginsdk.TypeString,
						Sensitive: true,
						Computed:  true,
					},

					"thumbprint": {
						Type:      pluginsdk.TypeString,
						Sensitive: true,
						Computed:  true,
					},
				},
			},
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
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

			codeSigningAccountId, err := codesigningaccounts.ParseCodeSigningAccountID(model.TrustedSigningAccountId)
			if err != nil {
				return err
			}

			id := certificateprofiles.NewCertificateProfileID(codeSigningAccountId.SubscriptionId, codeSigningAccountId.ResourceGroupName, codeSigningAccountId.CodeSigningAccountName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &certificateprofiles.CertificateProfile{
				Properties: &certificateprofiles.CertificateProfileProperties{
					IdentityValidationId: model.IdentityValidationId,
					IncludeCity:          &model.IncludeCity,
					IncludeCountry:       &model.IncludeCountry,
					IncludePostalCode:    &model.IncludePostalCode,
					IncludeState:         &model.IncludeState,
					IncludeStreetAddress: &model.IncludeStreetAddress,
					ProfileType:          model.ProfileType,
				},
			}

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r TrustedSigningCertificateProfileResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CodeSigning.Client.CertificateProfiles

			id, err := certificateprofiles.ParseCertificateProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model TrustedSigningCertificateProfileModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if metadata.ResourceData.HasChange("identity_validation_id") {
				properties.Properties.IdentityValidationId = model.IdentityValidationId
			}

			if metadata.ResourceData.HasChange("include_city") {
				properties.Properties.IncludeCity = &model.IncludeCity
			}

			if metadata.ResourceData.HasChange("include_country") {
				properties.Properties.IncludeCountry = &model.IncludeCountry
			}

			if metadata.ResourceData.HasChange("include_postal_code") {
				properties.Properties.IncludePostalCode = &model.IncludePostalCode
			}

			if metadata.ResourceData.HasChange("include_state") {
				properties.Properties.IncludeState = &model.IncludeState
			}

			if metadata.ResourceData.HasChange("include_street_address") {
				properties.Properties.IncludeStreetAddress = &model.IncludeStreetAddress
			}

			if metadata.ResourceData.HasChange("profile_type") {
				properties.Properties.ProfileType = model.ProfileType
			}

			if err := client.CreateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
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

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := TrustedSigningCertificateProfileModel{
				Name:                    id.CertificateProfileName,
				TrustedSigningAccountId: codesigningaccounts.NewCodeSigningAccountID(id.SubscriptionId, id.ResourceGroupName, id.CodeSigningAccountName).ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.Certificates = flattenCertificateModelArray(properties.Certificates)
					state.IdentityValidationId = properties.IdentityValidationId
					state.IncludeCity = pointer.From(properties.IncludeCity)
					state.IncludeCountry = pointer.From(properties.IncludeCountry)
					state.IncludePostalCode = pointer.From(properties.IncludePostalCode)
					state.IncludeState = pointer.From(properties.IncludeState)
					state.IncludeStreetAddress = pointer.From(properties.IncludeStreetAddress)
					state.ProfileType = properties.ProfileType
					state.Status = pointer.From(properties.Status)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r TrustedSigningCertificateProfileResource) Delete() sdk.ResourceFunc {
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

func flattenCertificateModelArray(inputList *[]certificateprofiles.Certificate) []CertificateModel {
	outputList := make([]CertificateModel, 0)
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := CertificateModel{
			Revocation:       flattenRevocationModel(input.Revocation),
			CreatedDate:      pointer.From(input.CreatedDate),
			EnhancedKeyUsage: pointer.From(input.EnhancedKeyUsage),
			ExpiryDate:       pointer.From(input.ExpiryDate),
			SerialNumber:     pointer.From(input.SerialNumber),
			Status:           pointer.From(input.Status),
			SubjectName:      pointer.From(input.SubjectName),
			Thumbprint:       pointer.From(input.Thumbprint),
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenRevocationModel(input *certificateprofiles.Revocation) []RevocationModel {
	var outputList []RevocationModel
	if input == nil {
		return outputList
	}
	output := RevocationModel{}
	if input.EffectiveAt != nil {
		output.EffectiveAt = *input.EffectiveAt
	}

	if input.FailureReason != nil {
		output.FailureReason = *input.FailureReason
	}

	if input.Reason != nil {
		output.Reason = *input.Reason
	}

	if input.Remarks != nil {
		output.Remarks = *input.Remarks
	}

	if input.RequestedAt != nil {
		output.RequestedAt = *input.RequestedAt
	}

	if input.Status != nil {
		output.Status = *input.Status
	}

	return append(outputList, output)
}
