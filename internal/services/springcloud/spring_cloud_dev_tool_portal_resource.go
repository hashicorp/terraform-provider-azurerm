// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

type SpringCloudDevToolPortalModel struct {
	Name                          string     `tfschema:"name"`
	SpringCloudServiceId          string     `tfschema:"spring_cloud_service_id"`
	ApplicationAcceleratorEnabled bool       `tfschema:"application_accelerator_enabled"`
	ApplicationLiveViewEnabled    bool       `tfschema:"application_live_view_enabled"`
	PublicNetworkAccessEnabled    bool       `tfschema:"public_network_access_enabled"`
	Sso                           []SsoModel `tfschema:"sso"`
}

type SsoModel struct {
	ClientId     string   `tfschema:"client_id"`
	ClientSecret string   `tfschema:"client_secret"`
	MetadataUrl  string   `tfschema:"metadata_url"`
	Scope        []string `tfschema:"scope"`
}

type SpringCloudDevToolPortalResource struct{}

func (s SpringCloudDevToolPortalResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_dev_tool_portal` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudDevToolPortalResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudDevToolPortalResource{}
)

func (s SpringCloudDevToolPortalResource) ResourceType() string {
	return "azurerm_spring_cloud_dev_tool_portal"
}

func (s SpringCloudDevToolPortalResource) ModelObject() interface{} {
	return &SpringCloudDevToolPortalModel{}
}

func (s SpringCloudDevToolPortalResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SpringCloudDevToolPortalID
}

func (s SpringCloudDevToolPortalResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"default"}, false),
		},

		"spring_cloud_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudServiceID,
		},

		"application_accelerator_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"application_live_view_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"sso": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"client_secret": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"metadata_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"scope": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func (s SpringCloudDevToolPortalResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudDevToolPortalResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudDevToolPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.DevToolPortalClient
			springId, err := parse.SpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := parse.NewSpringCloudDevToolPortalID(springId.SubscriptionId, springId.ResourceGroup, springId.SpringName, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			DevToolPortalResource := appplatform.DevToolPortalResource{
				Properties: &appplatform.DevToolPortalProperties{
					Public:        utils.Bool(model.PublicNetworkAccessEnabled),
					SsoProperties: expandSpringCloudDevToolPortalSsoProperties(model.Sso),
					Features:      expandSpringCloudDevToolPortalFeatures(model.ApplicationAcceleratorEnabled, model.ApplicationLiveViewEnabled),
				},
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName, DevToolPortalResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudDevToolPortalResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.DevToolPortalClient

			id, err := parse.SpringCloudDevToolPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudDevToolPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if _, err = client.Get(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName); err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			DevToolPortalResource := appplatform.DevToolPortalResource{
				Properties: &appplatform.DevToolPortalProperties{
					Public:        utils.Bool(model.PublicNetworkAccessEnabled),
					SsoProperties: expandSpringCloudDevToolPortalSsoProperties(model.Sso),
					Features:      expandSpringCloudDevToolPortalFeatures(model.ApplicationAcceleratorEnabled, model.ApplicationLiveViewEnabled),
				},
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName, DevToolPortalResource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s SpringCloudDevToolPortalResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.DevToolPortalClient

			id, err := parse.SpringCloudDevToolPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			state := SpringCloudDevToolPortalModel{
				Name:                 id.DevToolPortalName,
				SpringCloudServiceId: parse.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName).ID(),
			}

			var model SpringCloudDevToolPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if props := resp.Properties; props != nil {
				if props.Public != nil {
					state.PublicNetworkAccessEnabled = *props.Public
				}
				if props.SsoProperties != nil {
					state.Sso = flattenSpringCloudDevToolPortalSsoProperties(props.SsoProperties, model)
				}
				if props.Features != nil {
					if props.Features.ApplicationAccelerator != nil && props.Features.ApplicationAccelerator.State == appplatform.DevToolPortalFeatureStateEnabled {
						state.ApplicationAcceleratorEnabled = true
					}
					if props.Features.ApplicationLiveView != nil && props.Features.ApplicationLiveView.State == appplatform.DevToolPortalFeatureStateEnabled {
						state.ApplicationLiveViewEnabled = true
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudDevToolPortalResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.DevToolPortalClient

			id, err := parse.SpringCloudDevToolPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				if !response.WasNotFound(future.Response()) {
					return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func expandSpringCloudDevToolPortalFeatures(applicationAcceleratorEnabled, applicationLiveViewEnabled bool) *appplatform.DevToolPortalFeatureSettings {
	applicationAcceleratorState := appplatform.DevToolPortalFeatureStateDisabled
	if applicationAcceleratorEnabled {
		applicationAcceleratorState = appplatform.DevToolPortalFeatureStateEnabled
	}

	applicationLiveViewState := appplatform.DevToolPortalFeatureStateDisabled
	if applicationLiveViewEnabled {
		applicationLiveViewState = appplatform.DevToolPortalFeatureStateEnabled
	}

	return &appplatform.DevToolPortalFeatureSettings{
		ApplicationAccelerator: &appplatform.DevToolPortalFeatureDetail{
			State: applicationAcceleratorState,
		},
		ApplicationLiveView: &appplatform.DevToolPortalFeatureDetail{
			State: applicationLiveViewState,
		},
	}
}

func expandSpringCloudDevToolPortalSsoProperties(input []SsoModel) *appplatform.DevToolPortalSsoProperties {
	if len(input) == 0 {
		return nil
	}
	sso := input[0]

	return &appplatform.DevToolPortalSsoProperties{
		Scopes:       &sso.Scope,
		ClientID:     utils.String(sso.ClientId),
		ClientSecret: utils.String(sso.ClientSecret),
		MetadataURL:  utils.String(sso.MetadataUrl),
	}
}

func flattenSpringCloudDevToolPortalSsoProperties(properties *appplatform.DevToolPortalSsoProperties, model SpringCloudDevToolPortalModel) []SsoModel {
	if properties == nil {
		return []SsoModel{}
	}

	clientId := ""
	if properties.ClientID != nil {
		clientId = *properties.ClientID
	}

	clientSecret := ""
	if len(model.Sso) != 0 {
		clientSecret = model.Sso[0].ClientSecret
	}

	metadataUrl := ""
	if properties.MetadataURL != nil {
		metadataUrl = *properties.MetadataURL
	}

	scopes := make([]string, 0)
	if properties.Scopes != nil {
		scopes = *properties.Scopes
	}

	return []SsoModel{
		{
			ClientId:     clientId,
			ClientSecret: clientSecret,
			MetadataUrl:  metadataUrl,
			Scope:        scopes,
		},
	}
}
