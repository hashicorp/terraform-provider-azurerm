// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SpringCloudAPIPortalModel struct {
	Name                       string              `tfschema:"name"`
	SpringCloudServiceId       string              `tfschema:"spring_cloud_service_id"`
	GatewayIds                 []string            `tfschema:"gateway_ids"`
	HttpsOnlyEnabled           bool                `tfschema:"https_only_enabled"`
	InstanceCount              int64               `tfschema:"instance_count"`
	PublicNetworkAccessEnabled bool                `tfschema:"public_network_access_enabled"`
	ApiTryOutEnabled           bool                `tfschema:"api_try_out_enabled"`
	Sso                        []ApiPortalSsoModel `tfschema:"sso"`
	Url                        string              `tfschema:"url"`
}

type ApiPortalSsoModel struct {
	ClientId     string   `tfschema:"client_id"`
	ClientSecret string   `tfschema:"client_secret"`
	IssuerUri    string   `tfschema:"issuer_uri"`
	Scope        []string `tfschema:"scope"`
}

type SpringCloudAPIPortalResource struct{}

func (s SpringCloudAPIPortalResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_api_portal` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudAPIPortalResource{}
	_ sdk.ResourceWithStateMigration              = SpringCloudAPIPortalResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudAPIPortalResource{}
)

func (s SpringCloudAPIPortalResource) ResourceType() string {
	return "azurerm_spring_cloud_api_portal"
}

func (s SpringCloudAPIPortalResource) ModelObject() interface{} {
	return &SpringCloudAPIPortalModel{}
}

func (s SpringCloudAPIPortalResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateApiPortalID
}

func (s SpringCloudAPIPortalResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"default",
			}, false),
		},

		"spring_cloud_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSpringCloudServiceID,
		},

		"api_try_out_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"gateway_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: appplatform.ValidateGatewayID,
			},
		},

		"https_only_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"instance_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 500),
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
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"client_secret": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"issuer_uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"scope": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (s SpringCloudAPIPortalResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (s SpringCloudAPIPortalResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudApiPortalV0ToV1{},
		},
	}
}

func (s SpringCloudAPIPortalResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudAPIPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := appplatform.NewApiPortalID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			existing, err := client.ApiPortalsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			service, err := client.ServicesGet(ctx, *springId)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", springId, err)
			}
			if service.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", springId)
			}
			if service.Model.Sku == nil || service.Model.Sku.Name == nil || service.Model.Sku.Tier == nil {
				return fmt.Errorf("invalid `sku` for %s", springId)
			}

			apiTryOutEnabledState := appplatform.ApiPortalApiTryOutEnabledStateDisabled
			if model.ApiTryOutEnabled {
				apiTryOutEnabledState = appplatform.ApiPortalApiTryOutEnabledStateEnabled
			}

			apiPortalResource := appplatform.ApiPortalResource{
				Properties: &appplatform.ApiPortalProperties{
					GatewayIds:            pointer.To(model.GatewayIds),
					HTTPSOnly:             pointer.To(model.HttpsOnlyEnabled),
					Public:                pointer.To(model.PublicNetworkAccessEnabled),
					SsoProperties:         expandAPIPortalSsoProperties(model.Sso),
					ApiTryOutEnabledState: pointer.To(apiTryOutEnabledState),
				},
				Sku: &appplatform.Sku{
					Name:     service.Model.Sku.Name,
					Tier:     service.Model.Sku.Tier,
					Capacity: pointer.To(model.InstanceCount),
				},
			}
			err = client.ApiPortalsCreateOrUpdateThenPoll(ctx, id, apiPortalResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudAPIPortalResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudAPIPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApiPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ApiPortalsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			properties := resp.Model.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			sku := resp.Model.Sku
			if sku == nil {
				return fmt.Errorf("retrieving %s: sku was nil", id)
			}

			if metadata.ResourceData.HasChange("gateway_ids") {
				properties.GatewayIds = pointer.To(model.GatewayIds)
			}

			if metadata.ResourceData.HasChange("https_only_enabled") {
				properties.HTTPSOnly = pointer.To(model.HttpsOnlyEnabled)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				properties.Public = pointer.To(model.PublicNetworkAccessEnabled)
			}

			if metadata.ResourceData.HasChange("sso") {
				properties.SsoProperties = expandAPIPortalSsoProperties(model.Sso)
			}

			if metadata.ResourceData.HasChange("instance_count") {
				sku.Capacity = pointer.To(model.InstanceCount)
			}

			if metadata.ResourceData.HasChange("api_try_out_enabled") {
				apiTryOutEnabledState := appplatform.ApiPortalApiTryOutEnabledStateDisabled
				if model.ApiTryOutEnabled {
					apiTryOutEnabledState = appplatform.ApiPortalApiTryOutEnabledStateEnabled
				}
				properties.ApiTryOutEnabledState = pointer.To(apiTryOutEnabledState)
			}

			apiPortalResource := appplatform.ApiPortalResource{
				Properties: properties,
				Sku:        sku,
			}
			err = client.ApiPortalsCreateOrUpdateThenPoll(ctx, *id, apiPortalResource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s SpringCloudAPIPortalResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApiPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ApiPortalsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)

			var model SpringCloudAPIPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudAPIPortalModel{
				Name:                 id.ApiPortalName,
				SpringCloudServiceId: springId.ID(),
			}
			if resp.Model != nil {
				if props := resp.Model.Properties; props != nil {
					state.GatewayIds = flattenSpringCloudAPIPortalGatewayIds(props.GatewayIds)
					state.HttpsOnlyEnabled = pointer.From(props.HTTPSOnly)
					state.PublicNetworkAccessEnabled = pointer.From(props.Public)
					state.Sso = flattenAPIPortalSsoProperties(props.SsoProperties, model.Sso)
					state.ApiTryOutEnabled = props.ApiTryOutEnabledState != nil && *props.ApiTryOutEnabledState == appplatform.ApiPortalApiTryOutEnabledStateEnabled
				}

				if sku := resp.Model.Sku; sku != nil {
					state.InstanceCount = pointer.From(sku.Capacity)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudAPIPortalResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApiPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.ApiPortalsDeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandAPIPortalSsoProperties(input []ApiPortalSsoModel) *appplatform.SsoProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.SsoProperties{
		Scope:        pointer.To(v.Scope),
		ClientId:     pointer.To(v.ClientId),
		ClientSecret: pointer.To(v.ClientSecret),
		IssuerUri:    pointer.To(v.IssuerUri),
	}
}

func flattenAPIPortalSsoProperties(input *appplatform.SsoProperties, old []ApiPortalSsoModel) []ApiPortalSsoModel {
	if input == nil {
		return make([]ApiPortalSsoModel, 0)
	}

	oldItems := make(map[string]ApiPortalSsoModel)
	for _, item := range old {
		if item.IssuerUri != "" {
			oldItems[item.IssuerUri] = item
		}
	}

	issuerUri := pointer.From(input.IssuerUri)
	var clientId string
	var clientSecret string
	if oldItem, ok := oldItems[issuerUri]; ok {
		if oldItem.ClientId != "" {
			clientId = oldItem.ClientId
		}
		if oldItem.ClientSecret != "" {
			clientSecret = oldItem.ClientSecret
		}
	}
	return []ApiPortalSsoModel{
		{
			ClientId:     clientId,
			ClientSecret: clientSecret,
			IssuerUri:    issuerUri,
			Scope:        pointer.From(input.Scope),
		},
	}
}

func flattenSpringCloudAPIPortalGatewayIds(ids *[]string) []string {
	if ids == nil || len(*ids) == 0 {
		return nil
	}
	out := make([]string, 0)
	for _, id := range *ids {
		gatewayId, err := appplatform.ParseGatewayIDInsensitively(id)
		if err == nil {
			out = append(out, gatewayId.ID())
		}
	}
	return out
}
