// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arcresourcebridge

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ArcResourceBridgeApplianceResource{}

type ArcResourceBridgeApplianceResource struct{}

type ApplianceModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Location          string                         `tfschema:"location"`
	Distro            appliances.Distro              `tfschema:"distro"`
	Identity          []identity.ModelSystemAssigned `tfschema:"identity"`
	Provider          appliances.Provider            `tfschema:"infrastructure_provider"`
	PublicKeyBase64   string                         `tfschema:"public_key_base64"`
	Tags              map[string]interface{}         `tfschema:"tags"`
}

func (r ArcResourceBridgeApplianceResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 260),
				validation.StringMatch(regexp.MustCompile(`[^+#%&'?/,%\\]+$`), "any of '+', '#', '%', '&', ''', '?', '/', ',', '%', '&', '\\', are not allowed"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"distro": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(appliances.DistroAKSEdge),
			}, false),
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"infrastructure_provider": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(appliances.ProviderHCI),
				string(appliances.ProviderSCVMM),
				string(appliances.ProviderVMWare),
			}, false),
		},

		"public_key_base64": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.Base64EncodedString,
		},

		"tags": commonschema.Tags(),
	}
}

func (r ArcResourceBridgeApplianceResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ArcResourceBridgeApplianceResource) ModelObject() interface{} {
	return &ArcResourceBridgeApplianceResource{}
}

func (r ArcResourceBridgeApplianceResource) ResourceType() string {
	return "azurerm_arc_resource_bridge_appliance"
}

func (r ArcResourceBridgeApplianceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApplianceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ArcResourceBridge.AppliancesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := appliances.NewApplianceID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
			if err != nil {
				return err
			}

			parameters := appliances.Appliance{
				Location: location.Normalize(model.Location),
				Properties: &appliances.ApplianceProperties{
					Distro: pointer.To(model.Distro),
					InfrastructureConfig: &appliances.AppliancePropertiesInfrastructureConfig{
						Provider: pointer.To(model.Provider),
					},
				},
				Tags: tags.Expand(model.Tags),
			}

			parameters.Identity = identity

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// since the public key could not be set during creation, update after creation
			if model.PublicKeyBase64 != "" {
				parameters.Properties.PublicKey = pointer.To(model.PublicKeyBase64)

				if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ArcResourceBridgeApplianceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcResourceBridge.AppliancesClient

			id, err := appliances.ParseApplianceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", *id)
			var model ApplianceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("identity") {
				identity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
				if err != nil {
					return err
				}

				parameters.Identity = identity
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if metadata.ResourceData.HasChanges("public_key_base64") {
				if parameters.Properties == nil {
					parameters.Properties = &appliances.ApplianceProperties{}
				}
				parameters.Properties.PublicKey = pointer.To(model.PublicKeyBase64)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ArcResourceBridgeApplianceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := appliances.ParseApplianceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.ArcResourceBridge.AppliancesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %q: %+v", *id, err)
			}

			state := ApplianceModel{
				Name:              id.ApplianceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Identity = identity.FlattenSystemAssignedToModel(model.Identity)
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.Distro = pointer.From(props.Distro)
					state.PublicKeyBase64 = pointer.From(props.PublicKey)

					if infraConfig := props.InfrastructureConfig; infraConfig != nil {
						state.Provider = pointer.From(infraConfig.Provider)
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r ArcResourceBridgeApplianceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcResourceBridge.AppliancesClient
			id, err := appliances.ParseApplianceID(metadata.ResourceData.Id())
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

func (r ArcResourceBridgeApplianceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appliances.ValidateApplianceID
}
