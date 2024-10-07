// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	kvValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

type botBaseResource struct{}

func (br botBaseResource) arguments(fields map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	output := map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(botservice.SkuNameF0),
				string(botservice.SkuNameS1),
			}, false),
		},

		"microsoft_app_id": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"endpoint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"developer_app_insights_api_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_application_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cmk_key_vault_key_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: kvValidate.NestedItemIdWithOptionalVersion,
		},

		"microsoft_app_msi_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"microsoft_app_tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"microsoft_app_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(botservice.MsaAppTypeMultiTenant),
				string(botservice.MsaAppTypeSingleTenant),
				string(botservice.MsaAppTypeUserAssignedMSI),
			}, false),
		},

		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"luis_app_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"luis_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"streaming_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"icon_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "https://docs.botframework.com/static/devportal/client/images/bot-framework-default.png",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": tags.Schema(),
	}

	for k, v := range fields {
		output[k] = v
	}

	return output
}

func (br botBaseResource) attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (br botBaseResource) createFunc(resourceName, botKind string) sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewBotServiceID(subscriptionId, metadata.ResourceData.Get("resource_group_name").(string), metadata.ResourceData.Get("name").(string))

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return tf.ImportAsExistsError(resourceName, id.ID())
			}

			displayName := metadata.ResourceData.Get("display_name").(string)
			if displayName == "" {
				displayName = id.Name
			}

			publicNetworkEnabled := botservice.PublicNetworkAccessEnabled
			if !metadata.ResourceData.Get("public_network_access_enabled").(bool) {
				publicNetworkEnabled = botservice.PublicNetworkAccessDisabled
			}

			props := botservice.Bot{
				Location: utils.String(metadata.ResourceData.Get("location").(string)),
				Sku: &botservice.Sku{
					Name: botservice.SkuName(metadata.ResourceData.Get("sku").(string)),
				},
				Kind: botservice.Kind(botKind),
				Properties: &botservice.BotProperties{
					DisplayName:                       pointer.To(displayName),
					Endpoint:                          pointer.To(metadata.ResourceData.Get("endpoint").(string)),
					MsaAppID:                          pointer.To(metadata.ResourceData.Get("microsoft_app_id").(string)),
					DeveloperAppInsightKey:            pointer.To(metadata.ResourceData.Get("developer_app_insights_key").(string)),
					DeveloperAppInsightsAPIKey:        pointer.To(metadata.ResourceData.Get("developer_app_insights_api_key").(string)),
					DeveloperAppInsightsApplicationID: pointer.To(metadata.ResourceData.Get("developer_app_insights_application_id").(string)),
					DisableLocalAuth:                  pointer.To(!metadata.ResourceData.Get("local_authentication_enabled").(bool)),
					IsCmekEnabled:                     utils.Bool(false),
					CmekKeyVaultURL:                   pointer.To(metadata.ResourceData.Get("cmk_key_vault_key_url").(string)),
					LuisAppIds:                        utils.ExpandStringSlice(metadata.ResourceData.Get("luis_app_ids").([]interface{})),
					LuisKey:                           pointer.To(metadata.ResourceData.Get("luis_key").(string)),
					PublicNetworkAccess:               publicNetworkEnabled,
					IsStreamingSupported:              pointer.To(metadata.ResourceData.Get("streaming_endpoint_enabled").(bool)),
					IconURL:                           pointer.To(metadata.ResourceData.Get("icon_url").(string)),
				},
				Tags: tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{})),
			}

			if _, ok := metadata.ResourceData.GetOk("cmk_key_vault_key_url"); ok {
				props.Properties.IsCmekEnabled = utils.Bool(true)
			}

			if v, ok := metadata.ResourceData.GetOk("microsoft_app_type"); ok {
				props.Properties.MsaAppType = botservice.MsaAppType(v.(string))
			}

			if v, ok := metadata.ResourceData.GetOk("microsoft_app_tenant_id"); ok {
				props.Properties.MsaAppTenantID = pointer.To(v.(string))
			}

			if v, ok := metadata.ResourceData.GetOk("microsoft_app_msi_id"); ok {
				props.Properties.MsaAppMSIResourceID = pointer.To(v.(string))
			}

			if _, err := client.Create(ctx, id.ResourceGroup, id.Name, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (br botBaseResource) updateFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient
			id, err := parse.BotServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("display_name") {
				existing.Properties.DisplayName = utils.String(metadata.ResourceData.Get("display_name").(string))
			}

			if metadata.ResourceData.HasChange("endpoint") {
				existing.Properties.Endpoint = utils.String(metadata.ResourceData.Get("endpoint").(string))
			}

			if metadata.ResourceData.HasChange("developer_app_insights_key") {
				existing.Properties.DeveloperAppInsightKey = utils.String(metadata.ResourceData.Get("developer_app_insights_key").(string))
			}

			if metadata.ResourceData.HasChange("developer_app_insights_api_key") {
				existing.Properties.DeveloperAppInsightsAPIKey = utils.String(metadata.ResourceData.Get("developer_app_insights_api_key").(string))
			}

			if metadata.ResourceData.HasChange("developer_app_insights_application_id") {
				existing.Properties.DeveloperAppInsightsApplicationID = utils.String(metadata.ResourceData.Get("developer_app_insights_application_id").(string))
			}

			if metadata.ResourceData.HasChange("local_authentication_enabled") {
				existing.Properties.DisableLocalAuth = utils.Bool(!metadata.ResourceData.Get("local_authentication_enabled").(bool))
			}

			if metadata.ResourceData.HasChange("luis_app_ids") {
				existing.Properties.LuisAppIds = utils.ExpandStringSlice(metadata.ResourceData.Get("luis_app_ids").([]interface{}))
			}

			if metadata.ResourceData.HasChange("luis_key") {
				existing.Properties.LuisKey = utils.String(metadata.ResourceData.Get("luis_key").(string))
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if metadata.ResourceData.Get("public_network_access_enabled").(bool) {
					existing.Properties.PublicNetworkAccess = botservice.PublicNetworkAccessEnabled
				} else {
					existing.Properties.PublicNetworkAccess = botservice.PublicNetworkAccessDisabled
				}
			}

			if metadata.ResourceData.HasChange("streaming_endpoint_enabled") {
				existing.Properties.IsStreamingSupported = utils.Bool(metadata.ResourceData.Get("streaming_endpoint_enabled").(bool))
			}

			if metadata.ResourceData.HasChange("icon_url") {
				existing.Properties.IconURL = utils.String(metadata.ResourceData.Get("icon_url").(string))
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{}))
			}

			if _, err := client.Update(ctx, id.ResourceGroup, id.Name, existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (br botBaseResource) readFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient

			id, err := parse.BotServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.Name)
			metadata.ResourceData.Set("resource_group_name", id.ResourceGroup)
			metadata.ResourceData.Set("location", location.NormalizeNilable(resp.Location))

			sku := ""
			if v := resp.Sku; v != nil {
				sku = string(v.Name)
			}
			metadata.ResourceData.Set("sku", sku)

			metadata.ResourceData.Set("tags", tags.ToTypedObject(resp.Tags))

			// The API doesn't return this property, so we need to set the value from config into state
			if apiKey, ok := metadata.ResourceData.GetOk("developer_app_insights_api_key"); ok && apiKey.(string) != "" {
				metadata.ResourceData.Set("developer_app_insights_api_key", apiKey.(string))
			}

			if props := resp.Properties; props != nil {
				msAppId := ""
				if v := props.MsaAppID; v != nil {
					msAppId = *v
				}
				metadata.ResourceData.Set("microsoft_app_id", msAppId)

				displayName := ""
				if v := props.DisplayName; v != nil {
					displayName = *v
				}
				metadata.ResourceData.Set("display_name", displayName)

				endpoint := ""
				if v := props.Endpoint; v != nil {
					endpoint = *v
				}
				metadata.ResourceData.Set("endpoint", endpoint)

				key := ""
				if v := props.DeveloperAppInsightKey; v != nil {
					key = *v
				}
				metadata.ResourceData.Set("developer_app_insights_key", key)

				appInsightsId := ""
				if v := props.DeveloperAppInsightsApplicationID; v != nil {
					appInsightsId = *v
				}
				metadata.ResourceData.Set("developer_app_insights_application_id", appInsightsId)

				msaAppType := ""
				if v := props.MsaAppType; v != "" {
					msaAppType = string(v)
				}
				metadata.ResourceData.Set("microsoft_app_type", msaAppType)

				msaAppTenantId := ""
				if v := props.MsaAppTenantID; v != nil {
					msaAppTenantId = *v
				}
				metadata.ResourceData.Set("microsoft_app_tenant_id", msaAppTenantId)

				msaAppMSIId := ""
				if v := props.MsaAppMSIResourceID; v != nil {
					msaAppMSIId = *v
				}
				metadata.ResourceData.Set("microsoft_app_msi_id", msaAppMSIId)

				localAuthEnabled := true
				if v := props.DisableLocalAuth; v != nil {
					localAuthEnabled = !*v
				}
				metadata.ResourceData.Set("local_authentication_enabled", localAuthEnabled)

				publicNetworkAccessEnabled := true
				if v := props.PublicNetworkAccess; v != botservice.PublicNetworkAccessEnabled {
					publicNetworkAccessEnabled = false
				}
				metadata.ResourceData.Set("public_network_access_enabled", publicNetworkAccessEnabled)

				var luisAppIds []string
				if v := props.LuisAppIds; v != nil {
					luisAppIds = *v
				}
				metadata.ResourceData.Set("luis_app_ids", utils.FlattenStringSlice(&luisAppIds))

				streamingEndpointEnabled := false
				if v := props.IsStreamingSupported; v != nil {
					streamingEndpointEnabled = *v
				}
				metadata.ResourceData.Set("streaming_endpoint_enabled", streamingEndpointEnabled)

				metadata.ResourceData.Set("icon_url", pointer.From(props.IconURL))

				metadata.ResourceData.Set("cmk_key_vault_key_url", pointer.From(props.CmekKeyVaultURL))
			}

			return nil
		},
	}
}

func (br botBaseResource) deleteFunc() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Bot.BotClient
			id, err := parse.BotServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (br botBaseResource) importerFunc(expectKind string) sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Bot.BotClient

		id, err := parse.BotServiceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if actualKind := string(resp.Kind); actualKind != expectKind {
			return fmt.Errorf("bot has mismatched type, expected: %q, got %q", expectKind, actualKind)
		}

		return nil
	}
}
