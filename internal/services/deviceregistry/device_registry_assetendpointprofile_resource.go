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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = AssetEndpointProfileResource{}

type AssetEndpointProfileResource struct{}

type AssetEndpointProfileResourceModel struct {
	Name                                          string            `tfschema:"name"`
	ResourceGroupName                             string            `tfschema:"resource_group_name"`
	Location                                      string            `tfschema:"location"`
	Tags                                          map[string]string `tfschema:"tags"`
	ExtendedLocationName                          string            `tfschema:"extended_location_name"`
	ExtendedLocationType                          string            `tfschema:"extended_location_type"`
	TargetAddress                                 string            `tfschema:"target_address"`
	EndpointProfileType                           string            `tfschema:"endpoint_profile_type"`
	DiscoveredAssetEndpointProfileRef             string            `tfschema:"discovered_asset_endpoint_profile_ref"`
	AdditionalConfiguration                       string            `tfschema:"additional_configuration"`
	AuthenticationMethod                          string            `tfschema:"authentication_method"`
	UsernamePasswordCredentialsUsernameSecretName string            `tfschema:"username_password_credentials_username_secret_name"`
	UsernamePasswordCredentialsPasswordSecretName string            `tfschema:"username_password_credentials_password_secret_name"`
	X509CredentialsCertificateSecretName          string            `tfschema:"x509_credentials_certificate_secret_name"`
}

func (AssetEndpointProfileResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"extended_location_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"extended_location_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
		"discovered_asset_endpoint_profile_ref": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"additional_configuration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"authentication_method": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(assetendpointprofiles.PossibleValuesForAuthenticationMethod(), false),
		},
		"username_password_credentials_username_secret_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"username_password_credentials_password_secret_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"x509_credentials_certificate_secret_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			client := metadata.Client.DeviceRegistry.AssetEndpointProfileClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config AssetEndpointProfileResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := assetendpointprofiles.NewAssetEndpointProfileID(subscriptionId, config.ResourceGroupName, config.Name)

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
					Name: config.ExtendedLocationName,
					Type: config.ExtendedLocationType,
				},
				Properties: &assetendpointprofiles.AssetEndpointProfileProperties{
					TargetAddress:       config.TargetAddress,
					EndpointProfileType: config.EndpointProfileType,
				},
			}

			if config.DiscoveredAssetEndpointProfileRef != "" {
				param.Properties.DiscoveredAssetEndpointProfileRef = pointer.To(config.DiscoveredAssetEndpointProfileRef)
			}

			if config.AdditionalConfiguration != "" {
				param.Properties.AdditionalConfiguration = pointer.To(config.AdditionalConfiguration)
			}

			populateAuthenticationProperties(&param, config)

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
			client := metadata.Client.DeviceRegistry.AssetEndpointProfileClient

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

			authenticationMethodChanged := metadata.ResourceData.HasChange("authentication_method")
			usernameSecretNameChanged := metadata.ResourceData.HasChange("username_password_credentials_username_secret_name")
			passwordSecretNameChanged := metadata.ResourceData.HasChange("username_password_credentials_password_secret_name")
			certificateSecretNameChanged := metadata.ResourceData.HasChange("x509_credentials_certificate_secret_name")
			if authenticationMethodChanged || usernameSecretNameChanged || passwordSecretNameChanged || certificateSecretNameChanged {
				authentication := assetendpointprofiles.AuthenticationUpdate{}
				param.Properties.Authentication = &authentication

				if authenticationMethodChanged {
					authentication.Method = pointer.To(assetendpointprofiles.AuthenticationMethod(config.AuthenticationMethod))
				}

				if usernameSecretNameChanged || passwordSecretNameChanged {
					usernamePasswordCreds := assetendpointprofiles.UsernamePasswordCredentialsUpdate{}

					if usernameSecretNameChanged {
						usernamePasswordCreds.UsernameSecretName = pointer.To(config.UsernamePasswordCredentialsUsernameSecretName)
					}

					if passwordSecretNameChanged {
						usernamePasswordCreds.PasswordSecretName = pointer.To(config.UsernamePasswordCredentialsPasswordSecretName)
					}

					_, usernameOk := metadata.ResourceData.GetOk("username_password_credentials_username_secret_name")
					_, passwordOk := metadata.ResourceData.GetOk("username_password_credentials_password_secret_name")
					if (usernameSecretNameChanged && usernameOk) || (passwordSecretNameChanged && passwordOk) {
						authentication.UsernamePasswordCredentials = &usernamePasswordCreds
					} else {
						authentication.UsernamePasswordCredentials = nil
					}
				}

				if certificateSecretNameChanged {
					if _, ok := metadata.ResourceData.GetOk("x509_credentials_certificate_secret_name"); ok {
						authentication.X509Credentials = &assetendpointprofiles.X509CredentialsUpdate{
							CertificateSecretName: pointer.To(config.X509CredentialsCertificateSecretName),
						}
					} else {
						authentication.X509Credentials = nil
					}
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
			client := metadata.Client.DeviceRegistry.AssetEndpointProfileClient

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

			// Convert the ARM model to the TF model
			state := AssetEndpointProfileResourceModel{
				Name:              id.AssetEndpointProfileName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.ExtendedLocationName = model.ExtendedLocation.Name
				state.ExtendedLocationType = model.ExtendedLocation.Type

				if props := model.Properties; props != nil {
					state.TargetAddress = props.TargetAddress
					state.EndpointProfileType = props.EndpointProfileType
					state.DiscoveredAssetEndpointProfileRef = pointer.From(props.DiscoveredAssetEndpointProfileRef)
					state.AdditionalConfiguration = pointer.From(props.AdditionalConfiguration)

					if auth := props.Authentication; auth != nil {
						state.AuthenticationMethod = string(auth.Method)

						if x509 := auth.X509Credentials; x509 != nil {
							state.X509CredentialsCertificateSecretName = x509.CertificateSecretName
						}

						if up := auth.UsernamePasswordCredentials; up != nil {
							state.UsernamePasswordCredentialsUsernameSecretName = up.UsernameSecretName
							state.UsernamePasswordCredentialsPasswordSecretName = up.PasswordSecretName
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
			client := metadata.Client.DeviceRegistry.AssetEndpointProfileClient

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

func populateAuthenticationProperties(param *assetendpointprofiles.AssetEndpointProfile, config AssetEndpointProfileResourceModel) {
	// If the authentication method is not set, we don't need to populate the authentication properties
	if config.AuthenticationMethod == "" {
		return
	}

	param.Properties.Authentication = &assetendpointprofiles.Authentication{
		Method: assetendpointprofiles.AuthenticationMethod(config.AuthenticationMethod),
	}

	if config.X509CredentialsCertificateSecretName != "" {
		param.Properties.Authentication.X509Credentials = &assetendpointprofiles.X509Credentials{
			CertificateSecretName: config.X509CredentialsCertificateSecretName,
		}
	}

	if config.UsernamePasswordCredentialsUsernameSecretName != "" || config.UsernamePasswordCredentialsPasswordSecretName != "" {
		param.Properties.Authentication.UsernamePasswordCredentials = &assetendpointprofiles.UsernamePasswordCredentials{
			UsernameSecretName: config.UsernamePasswordCredentialsUsernameSecretName,
			PasswordSecretName: config.UsernamePasswordCredentialsPasswordSecretName,
		}
	}
}
