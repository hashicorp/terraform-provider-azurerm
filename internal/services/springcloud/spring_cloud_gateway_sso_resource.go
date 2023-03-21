package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/appplatform/2022-11-01-preview/appplatform"
)

type SpringCloudGatewaySsoModel struct {
	SpringCloudGatewayId string   `tfschema:"spring_cloud_gateway_id"`
	Scope                []string `tfschema:"scope"`
	ClientID             string   `tfschema:"client_id"`
	ClientSecret         string   `tfschema:"client_secret"`
	IssuerUri            string   `tfschema:"issuer_uri"`
}

type SpringCloudGatewaySsoResource struct{}

var _ sdk.ResourceWithUpdate = SpringCloudGatewaySsoResource{}

func (s SpringCloudGatewaySsoResource) ResourceType() string {
	return "azurerm_spring_cloud_gateway_sso"
}

func (s SpringCloudGatewaySsoResource) ModelObject() interface{} {
	return &SpringCloudGatewaySsoModel{}
}

func (s SpringCloudGatewaySsoResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SpringCloudGatewayID
}

func (s SpringCloudGatewaySsoResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"spring_cloud_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudGatewayID,
		},

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

		"issuer_uri": {
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
	}
}

func (s SpringCloudGatewaySsoResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudGatewaySsoResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudGatewaySsoModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.GatewayClient
			id, err := parse.SpringCloudGatewayID(model.SpringCloudGatewayId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName)
			if err != nil {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if existing.Properties != nil && existing.Properties.SsoProperties != nil {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			existing.Properties.SsoProperties = &appplatform.SsoProperties{
				Scope:        utils.ToPtr(model.Scope),
				ClientID:     utils.String(model.ClientID),
				ClientSecret: utils.String(model.ClientSecret),
				IssuerURI:    utils.String(model.IssuerUri),
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, existing)
			if err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudGatewaySsoResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudGatewaySsoModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.GatewayClient
			id, err := parse.SpringCloudGatewayID(model.SpringCloudGatewayId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName)
			if err != nil {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if existing.Properties.SsoProperties == nil {
				existing.Properties.SsoProperties = &appplatform.SsoProperties{}
			}
			if metadata.ResourceData.HasChange("scope") {
				existing.Properties.SsoProperties.Scope = utils.ToPtr(model.Scope)
			}
			if metadata.ResourceData.HasChange("client_id") {
				existing.Properties.SsoProperties.ClientID = utils.String(model.ClientID)
			}
			if metadata.ResourceData.HasChange("client_secret") {
				existing.Properties.SsoProperties.ClientSecret = utils.String(model.ClientSecret)
			}
			if metadata.ResourceData.HasChange("issuer_uri") {
				existing.Properties.SsoProperties.IssuerURI = utils.String(model.IssuerUri)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, existing)
			if err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (s SpringCloudGatewaySsoResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudGatewaySsoModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.AppPlatform.GatewayClient

			id, err := parse.SpringCloudGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			state := SpringCloudGatewaySsoModel{
				SpringCloudGatewayId: model.SpringCloudGatewayId,
				ClientID:             model.ClientID,
				ClientSecret:         model.ClientSecret,
			}

			if props := resp.Properties; props != nil {
				if sso := resp.Properties.SsoProperties; sso != nil {
					if sso.Scope != nil {
						state.Scope = *sso.Scope
					}
					if sso.IssuerURI != nil {
						state.IssuerUri = *sso.IssuerURI
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudGatewaySsoResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.GatewayClient
			id, err := parse.SpringCloudGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName)
			if err != nil {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			existing.Properties.SsoProperties = nil

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, existing)
			if err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
			}

			return nil
		},
	}
}
