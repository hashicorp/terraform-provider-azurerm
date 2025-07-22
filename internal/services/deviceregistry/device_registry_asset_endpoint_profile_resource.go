package deviceregistry

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	resourceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	resourceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	AssetEndpointProfileExtendedLocationTypeCustomLocation = "CustomLocation"
)

var _ sdk.Resource = AssetEndpointProfileResource{}

type AssetEndpointProfileResource struct{}

type AssetEndpointProfileResourceModel struct {
	Name                                    string                `tfschema:"name"`
	ResourceGroupId                         string                `tfschema:"resource_group_id"`
	Location                                string                `tfschema:"location"`
	Tags                                    map[string]string     `tfschema:"tags"`
	ExtendedLocationId                      string                `tfschema:"extended_location_id"`
	TargetAddress                           string                `tfschema:"target_address"`
	EndpointProfileType                     string                `tfschema:"endpoint_profile_type"`
	DiscoveredAssetEndpointProfileReference string                `tfschema:"discovered_asset_endpoint_profile_reference"`
	AdditionalConfiguration                 string                `tfschema:"additional_configuration"`
	Authentication                          []AuthenticationModel `tfschema:"authentication"`
}

type AuthenticationModel struct {
	Method                                        string `tfschema:"method"`
	UsernamePasswordCredentialsUsernameSecretName string `tfschema:"username_password_credential_username_secret_name"`
	UsernamePasswordCredentialsPasswordSecretName string `tfschema:"username_password_credential_password_secret_name"`
	X509CredentialsCertificateSecretName          string `tfschema:"x509_credential_certificate_secret_name"`
}

func (AssetEndpointProfileResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: resourceValidate.ResourceGroupID,
		},
		"extended_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},
		"target_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"endpoint_profile_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"discovered_asset_endpoint_profile_reference": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"additional_configuration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"method": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(assetendpointprofiles.PossibleValuesForAuthenticationMethod(), false),
					},
					"username_password_credential_username_secret_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"username_password_credential_password_secret_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"x509_credential_certificate_secret_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
		"location": commonschema.Location(),
		"tags":     commonschema.Tags(),
	}
}

func (AssetEndpointProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AssetEndpointProfileResource) ModelObject() interface{} {
	return &AssetEndpointProfileResourceModel{}
}

func (AssetEndpointProfileResource) ResourceType() string {
	return "azurerm_device_registry_asset_endpoint_profile"
}

func (r AssetEndpointProfileResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DeviceRegistry.AssetEndpointProfilesClient

			var config AssetEndpointProfileResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resourceGroupId, err := resourceParse.ResourceGroupID(config.ResourceGroupId)
			if err != nil {
				return err
			}

			id := assetendpointprofiles.NewAssetEndpointProfileID(resourceGroupId.SubscriptionId, resourceGroupId.ResourceGroup, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// Convert the TF model to the ARM model
			// Optional ARM resource properties are pointers.
			param := assetendpointprofiles.AssetEndpointProfile{
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
				ExtendedLocation: assetendpointprofiles.ExtendedLocation{
					Name: config.ExtendedLocationId,
					Type: AssetEndpointProfileExtendedLocationTypeCustomLocation,
				},
				Properties: &assetendpointprofiles.AssetEndpointProfileProperties{
					TargetAddress:       config.TargetAddress,
					EndpointProfileType: config.EndpointProfileType,
				},
			}

			if config.DiscoveredAssetEndpointProfileReference != "" {
				param.Properties.DiscoveredAssetEndpointProfileRef = pointer.To(config.DiscoveredAssetEndpointProfileReference)
			}

			if config.AdditionalConfiguration != "" {
				param.Properties.AdditionalConfiguration = pointer.To(config.AdditionalConfiguration)
			}

			param.Properties.Authentication = expandAuthentication(config.Authentication)

			if err := client.CreateOrReplaceThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AssetEndpointProfileResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DeviceRegistry.AssetEndpointProfilesClient

			id, err := assetendpointprofiles.ParseAssetEndpointProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config AssetEndpointProfileResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Convert the TF model to the ARM model
			param := assetendpointprofiles.AssetEndpointProfileUpdate{
				Properties: &assetendpointprofiles.AssetEndpointProfileUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(config.Tags)
			}

			if metadata.ResourceData.HasChange("target_address") {
				param.Properties.TargetAddress = pointer.To(config.TargetAddress)
			}

			if metadata.ResourceData.HasChange("endpoint_profile_type") {
				param.Properties.EndpointProfileType = pointer.To(config.EndpointProfileType)
			}

			if metadata.ResourceData.HasChange("additional_configuration") {
				param.Properties.AdditionalConfiguration = pointer.To(config.AdditionalConfiguration)
			}

			if metadata.ResourceData.HasChange("authentication") {
				authenticationModel := config.Authentication[0]
				authentication := &assetendpointprofiles.AuthenticationUpdate{}
				param.Properties.Authentication = authentication

				if metadata.ResourceData.HasChange("authentication.0.method") {
					authentication.Method = pointer.To(assetendpointprofiles.AuthenticationMethod(authenticationModel.Method))
				}

				usernameSecretChanged := metadata.ResourceData.HasChange("authentication.0.username_password_credential_username_secret_name")
				passwordSecretChanged := metadata.ResourceData.HasChange("authentication.0.username_password_credential_password_secret_name")
				if usernameSecretChanged || passwordSecretChanged {
					usernamePasswordCreds := assetendpointprofiles.UsernamePasswordCredentialsUpdate{}
					authentication.UsernamePasswordCredentials = &usernamePasswordCreds

					if usernameSecretChanged {
						usernamePasswordCreds.UsernameSecretName = pointer.To(authenticationModel.UsernamePasswordCredentialsUsernameSecretName)
					}

					if passwordSecretChanged {
						usernamePasswordCreds.PasswordSecretName = pointer.To(authenticationModel.UsernamePasswordCredentialsPasswordSecretName)
					}
				}

				if metadata.ResourceData.HasChange("authentication.0.x509_credential_certificate_secret_name") {
					authentication.X509Credentials = &assetendpointprofiles.X509CredentialsUpdate{}
					authentication.X509Credentials.CertificateSecretName = pointer.To(authenticationModel.X509CredentialsCertificateSecretName)
				}
			}

			if err := client.UpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (AssetEndpointProfileResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DeviceRegistry.AssetEndpointProfilesClient

			id, err := assetendpointprofiles.ParseAssetEndpointProfileID(metadata.ResourceData.Id())
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

			resourceGroupId := resourceParse.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName)

			// Convert the ARM model to the TF model
			state := AssetEndpointProfileResourceModel{
				Name:            id.AssetEndpointProfileName,
				ResourceGroupId: resourceGroupId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.ExtendedLocationId = model.ExtendedLocation.Name

				if props := model.Properties; props != nil {
					state.TargetAddress = props.TargetAddress
					state.EndpointProfileType = props.EndpointProfileType
					state.DiscoveredAssetEndpointProfileReference = pointer.From(props.DiscoveredAssetEndpointProfileRef)
					state.AdditionalConfiguration = pointer.From(props.AdditionalConfiguration)

					if auth := props.Authentication; auth != nil {
						authenticationModel := AuthenticationModel{
							Method: string(auth.Method),
						}

						if usernamePassword := auth.UsernamePasswordCredentials; usernamePassword != nil {
							authenticationModel.UsernamePasswordCredentialsUsernameSecretName = usernamePassword.UsernameSecretName
							authenticationModel.UsernamePasswordCredentialsPasswordSecretName = usernamePassword.PasswordSecretName
						}

						if x509 := auth.X509Credentials; x509 != nil {
							authenticationModel.X509CredentialsCertificateSecretName = x509.CertificateSecretName
						}

						state.Authentication = []AuthenticationModel{
							authenticationModel,
						}
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (AssetEndpointProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DeviceRegistry.AssetEndpointProfilesClient

			id, err := assetendpointprofiles.ParseAssetEndpointProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (AssetEndpointProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return assetendpointprofiles.ValidateAssetEndpointProfileID
}

func expandAuthentication(authenticationModels []AuthenticationModel) *assetendpointprofiles.Authentication {
	if len(authenticationModels) == 0 {
		return nil
	}

	authenticationModel := authenticationModels[0]
	authentication := &assetendpointprofiles.Authentication{
		Method: assetendpointprofiles.AuthenticationMethod(authenticationModel.Method),
	}

	if authenticationModel.UsernamePasswordCredentialsUsernameSecretName != "" || authenticationModel.UsernamePasswordCredentialsPasswordSecretName != "" {
		authentication.UsernamePasswordCredentials = &assetendpointprofiles.UsernamePasswordCredentials{
			UsernameSecretName: authenticationModel.UsernamePasswordCredentialsUsernameSecretName,
			PasswordSecretName: authenticationModel.UsernamePasswordCredentialsPasswordSecretName,
		}
	}

	if authenticationModel.X509CredentialsCertificateSecretName != "" {
		authentication.X509Credentials = &assetendpointprofiles.X509Credentials{
			CertificateSecretName: authenticationModel.X509CredentialsCertificateSecretName,
		}
	}

	return authentication
}
