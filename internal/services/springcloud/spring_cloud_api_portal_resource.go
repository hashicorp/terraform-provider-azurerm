// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

type SpringCloudAPIPortalModel struct {
	Name                       string              `tfschema:"name"`
	SpringCloudServiceId       string              `tfschema:"spring_cloud_service_id"`
	GatewayIds                 []string            `tfschema:"gateway_ids"`
	HttpsOnlyEnabled           bool                `tfschema:"https_only_enabled"`
	InstanceCount              int                 `tfschema:"instance_count"`
	PublicNetworkAccessEnabled bool                `tfschema:"public_network_access_enabled"`
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

var _ sdk.ResourceWithUpdate = SpringCloudAPIPortalResource{}
var _ sdk.ResourceWithStateMigration = SpringCloudAPIPortalResource{}

func (s SpringCloudAPIPortalResource) ResourceType() string {
	return "azurerm_spring_cloud_api_portal"
}

func (s SpringCloudAPIPortalResource) ModelObject() interface{} {
	return &SpringCloudAPIPortalModel{}
}

func (s SpringCloudAPIPortalResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SpringCloudAPIPortalID
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
			ValidateFunc: validate.SpringCloudServiceID,
		},

		"gateway_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.SpringCloudGatewayID,
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

			client := metadata.Client.AppPlatform.APIPortalClient
			servicesClient := metadata.Client.AppPlatform.ServicesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			springId, err := parse.SpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := parse.NewSpringCloudAPIPortalID(subscriptionId, springId.ResourceGroup, springId.SpringName, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_api_portal", id.ID())
			}

			service, err := servicesClient.Get(ctx, springId.ResourceGroup, springId.SpringName)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", springId, err)
			}
			if service.Sku == nil || service.Sku.Name == nil || service.Sku.Tier == nil {
				return fmt.Errorf("invalid `sku` for %s", springId)
			}

			apiPortalResource := appplatform.APIPortalResource{
				Properties: &appplatform.APIPortalProperties{
					GatewayIds:    pointer.To(model.GatewayIds),
					HTTPSOnly:     pointer.To(model.HttpsOnlyEnabled),
					Public:        pointer.To(model.PublicNetworkAccessEnabled),
					SsoProperties: expandAPIPortalSsoProperties(model.Sso),
				},
				Sku: &appplatform.Sku{
					Name:     service.Sku.Name,
					Tier:     service.Sku.Tier,
					Capacity: pointer.To(int32(model.InstanceCount)),
				},
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName, apiPortalResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
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

			client := metadata.Client.AppPlatform.APIPortalClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			springId, err := parse.SpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := parse.NewSpringCloudAPIPortalID(subscriptionId, springId.ResourceGroup, springId.SpringName, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
			}
			if utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("retrieving %s: resource was not found", id)
			}

			if existing.Properties == nil {
				return fmt.Errorf("retrieving %s: properties are nil", id)
			}
			properties := existing.Properties

			if existing.Sku == nil {
				return fmt.Errorf("retrieving %s: sku is nil", id)
			}
			sku := existing.Sku

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
				sku.Capacity = pointer.To(int32(model.InstanceCount))
			}

			apiPortalResource := appplatform.APIPortalResource{
				Properties: properties,
				Sku:        sku,
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName, apiPortalResource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s SpringCloudAPIPortalResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.APIPortalClient

			id, err := parse.SpringCloudAPIPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					log.Printf("[INFO] %q does not exist - removing from state", id.ID())
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			springId := parse.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName)

			var model SpringCloudAPIPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudAPIPortalModel{
				Name:                 id.ApiPortalName,
				SpringCloudServiceId: springId.ID(),
			}
			if resp.Sku != nil {
				state.InstanceCount = int(pointer.From(resp.Sku.Capacity))
			}
			if props := resp.Properties; props != nil {
				state.GatewayIds = flattenSpringCloudAPIPortalGatewayIds(props.GatewayIds)
				state.HttpsOnlyEnabled = pointer.From(props.HTTPSOnly)
				state.PublicNetworkAccessEnabled = pointer.From(props.Public)
				state.Sso = flattenAPIPortalSsoProperties(props.SsoProperties, model.Sso)
				state.Url = pointer.From(props.URL)
			}
			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudAPIPortalResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.APIPortalClient

			id, err := parse.SpringCloudAPIPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
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
		ClientID:     pointer.To(v.ClientId),
		ClientSecret: pointer.To(v.ClientSecret),
		IssuerURI:    pointer.To(v.IssuerUri),
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

	var issuerUri string
	if input.IssuerURI != nil {
		issuerUri = *input.IssuerURI
	}
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
		gatewayId, err := parse.SpringCloudGatewayIDInsensitively(id)
		if err == nil {
			out = append(out, gatewayId.ID())
		}
	}
	return out
}
