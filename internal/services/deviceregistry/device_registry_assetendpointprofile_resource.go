package deviceregistry

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = AssetEndpointProfileResource{}

type AssetEndpointProfileResource struct{}

type AssetEndpointProfileResourceModel struct {
	Name                                          string                     `tfschema:"name"`
	ResourceGroupName                             string                     `tfschema:"resource_group_name"`
	Type                         string                   `tfschema:"type"`
	Location                                      string                     `tfschema:"location"`
	Tags                                          map[string]string          `tfschema:"tags"`
	ExtendedLocationName                          string                     `tfschema:"extended_location_name"`
	ExtendedLocationType                          string                     `tfschema:"extended_location_type"`
	ProvisioningState                             string                     `tfschema:"provisioning_state"`
	Uuid                                          string                     `tfschema:"uuid"`
	TargetAddress                                 string                     `tfschema:"target_address"`
	EndpointProfileType                           string                     `tfschema:"endpoint_profile_type"`
	DiscoveredAssetEndpointProfileRef             string                     `tfschema:"discovered_asset_endpoint_profile_ref"`
	AdditionalConfiguration                       string                     `tfschema:"additional_configuration"`
	AuthenticationMethod                          string                     `tfschema:"authentication_method"`
	UsernamePasswordCredentialsUsernameSecretName string                     `tfschema:"username_password_credentials_username_secret_name"`
	UsernamePasswordCredentialsPasswordSecretName string                     `tfschema:"username_password_credentials_password_secret_name"`
	X509CredentialsCertificateSecretName          string                     `tfschema:"x509_credentials_certificate_secret_name"`
	Status                                        AssetEndpointProfileStatus `tfschema:"status"`
}

type AssetEndpointProfileStatus struct {
	Errors []StatusError `tfschema:"errors"`
}

type StatusError struct {
	Code    int64  `tfschema:"code"`
	Message string `tfschema:"message"`
}

func (AssetEndpointProfileResource) Arguments() map[string]*pluginsdk.Schema {
	// add the other assetendpointprofile properties
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default: string(assetendpointprofiles.AuthenticationMethodCertificate),
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
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"uuid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"status": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: map[string]*pluginsdk.Schema{
				"errors": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeMap,
						Elem: map[string]*pluginsdk.Schema{
							"code": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"message": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
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
					TargetAddress:                     config.TargetAddress,
					EndpointProfileType:               config.EndpointProfileType,
					DiscoveredAssetEndpointProfileRef: pointer.To(config.DiscoveredAssetEndpointProfileRef),
					AdditionalConfiguration:           pointer.To(config.AdditionalConfiguration),
				},
			}
			populateAuthenticationProperties(&param, config)

			if _, err := client.CreateOrReplace(ctx, id, param); err != nil {
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

			id, err := assetendpointprofiles.ParseAssetEndpointProfileID(metadata.ResourceData.Get("id").(string))
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

			if metadata.ResourceData.HasChange("authentication_method") {
				param.Properties.Authentication = &assetendpointprofiles.Authentication{
					Method: pointer.To(config.AuthenticationMethod),
				}
			}

			if metadata.ResourceData.HasChange("username_password_credentials_username_secret_name") {
				if param.Properties.Authentication == nil {
					param.Properties.Authentication = &assetendpointprofiles.Authentication{}
				}
				param.Properties.Authentication.UsernamePasswordCredentials = &assetendpointprofiles.UsernamePasswordCredentials{
					UsernameSecretName: pointer.To(config.UsernamePasswordCredentialsUsernameSecretName),
				}
			}

			if metadata.ResourceData.HasChange("username_password_credentials_password_secret_name") {
				if param.Properties.Authentication == nil {
					param.Properties.Authentication = &assetendpointprofiles.Authentication{}
				}
				if param.Properties.Authentication.UsernamePasswordCredentials == nil {
					param.Properties.Authentication.UsernamePasswordCredentials = &assetendpointprofiles.UsernamePasswordCredentials{}
				}
				param.Properties.Authentication.UsernamePasswordCredentials.PasswordSecretName = pointer.To(config.UsernamePasswordCredentialsPasswordSecretName)
			}

			if metadata.ResourceData.HasChange("x509_credentials_certificate_secret_name") {
				if param.Properties.Authentication == nil {
					param.Properties.Authentication = &assetendpointprofiles.Authentication{}
				}
				param.Properties.Authentication.X509Credentials = &assetendpointprofiles.X509Credentials{
					CertificateSecretName: pointer.To(config.X509CredentialsCertificateSecretName),
				}
			}

			if _, err := client.Update(ctx, *id, param); err != nil {
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
				state.Location = location.NormalizeNilable(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Type = pointer.From(model.Type)
				if props := model.Properties; props != nil {
					state.Uuid = pointer.From(props.Uuid)
					state.ExtendedLocationName = pointer.From(props.ExtendedLocation.Name)
					state.ExtendedLocationType = pointer.From(props.ExtendedLocation.Type)
					state.ProvisioningState = pointer.From(props.ProvisioningState)
					state.TargetAddress = props.TargetAddress
					state.EndpointProfileType = props.EndpointProfileType
					state.DiscoveredAssetEndpointProfileRef = pointer.From(props.DiscoveredAssetEndpointProfileRef)
					state.AdditionalConfiguration = pointer.From(props.AdditionalConfiguration)
					if auth := props.Authentication; auth != nil {
						state.AuthenticationMethod = auth.Method
						if x509 := auth.X509Credentials; x509 != nil {
							state.X509CredentialsCertificateSecretName = x509.CertificateSecretName
						}
						if up := auth.UsernamePasswordCredentials; up != nil {
							state.UsernamePasswordCredentialsUsernameSecretName = up.UsernameSecretName
							state.UsernamePasswordCredentialsPasswordSecretName = up.PasswordSecretName
						}
					}
					if status := model.Status; status != nil {
						state.Status = AssetEndpointProfileStatus{
							Errors: toTFAssetEndpointProfileErrorStatuses(status.Errors),
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
	return assetendpointprofiles.ValidateAssetID
}

func populateAuthenticationProperties(param *assetendpointprofiles.AssetEndpointProfile, config AssetEndpointProfileResourceModel) {
	param.Properties.Authentication = &assetendpointprofiles.Authentication{
		Method: config.AuthenticationMethod,
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

func toTFAssetEndpointProfileErrorStatuses(errors []assetendpointprofiles.StatusError) []StatusError {
	var tfErrors []StatusError
	for _, err := range errors {
		tfErrors = append(tfErrors, StatusError{
			Code:    err.Code,
			Message: err.Message,
		})
	}
	return tfErrors
}