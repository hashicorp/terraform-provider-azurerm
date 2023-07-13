// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type IotHubDeviceUpdateAccountResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotHubDeviceUpdateAccountResource{}
)

type IotHubDeviceUpdateAccountModel struct {
	Name                       string            `tfschema:"name"`
	ResourceGroupName          string            `tfschema:"resource_group_name"`
	Location                   string            `tfschema:"location"`
	HostName                   string            `tfschema:"host_name"`
	PublicNetworkAccessEnabled bool              `tfschema:"public_network_access_enabled"`
	Sku                        deviceupdates.SKU `tfschema:"sku"`
	Tags                       map[string]string `tfschema:"tags"`
}

func (r IotHubDeviceUpdateAccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubDeviceUpdateAccountName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  deviceupdates.SKUStandard,
			ValidateFunc: validation.StringInSlice([]string{
				string(deviceupdates.SKUFree),
				string(deviceupdates.SKUStandard),
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r IotHubDeviceUpdateAccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r IotHubDeviceUpdateAccountResource) ResourceType() string {
	return "azurerm_iothub_device_update_account"
}

func (r IotHubDeviceUpdateAccountResource) ModelObject() interface{} {
	return &IotHubDeviceUpdateAccountModel{}
}

func (r IotHubDeviceUpdateAccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deviceupdates.ValidateAccountID
}

func (r IotHubDeviceUpdateAccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model IotHubDeviceUpdateAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.IoTHub.DeviceUpdatesClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := deviceupdates.NewAccountID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.AccountsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			publicNetworkAccess := deviceupdates.PublicNetworkAccessEnabled
			if !model.PublicNetworkAccessEnabled {
				publicNetworkAccess = deviceupdates.PublicNetworkAccessDisabled
			}

			input := &deviceupdates.Account{
				Location: location.Normalize(model.Location),
				Identity: identityValue,
				Properties: &deviceupdates.AccountProperties{
					PublicNetworkAccess: &publicNetworkAccess,
					Sku:                 &model.Sku,
				},
				Tags: &model.Tags,
			}

			if err := client.AccountsCreateThenPoll(ctx, id, *input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r IotHubDeviceUpdateAccountResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := IotHubDeviceUpdateAccountModel{
				Name:              id.AccountName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if properties := model.Properties; properties != nil {
				if properties.HostName != nil {
					state.HostName = *properties.HostName
				}

				publicNetworkAccessEnabled := true
				if properties.PublicNetworkAccess != nil && *properties.PublicNetworkAccess == deviceupdates.PublicNetworkAccessDisabled {
					publicNetworkAccessEnabled = false
				}
				state.PublicNetworkAccessEnabled = publicNetworkAccessEnabled

				sku := deviceupdates.SKUStandard
				if properties.Sku != nil {
					sku = *properties.Sku
				}
				state.Sku = sku
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r IotHubDeviceUpdateAccountResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model IotHubDeviceUpdateAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.AccountsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			existing := resp.Model
			if existing == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if existing.Properties == nil {
				existing.Properties = &deviceupdates.AccountProperties{}
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				publicNetworkAccess := deviceupdates.PublicNetworkAccessEnabled
				if !model.PublicNetworkAccessEnabled {
					publicNetworkAccess = deviceupdates.PublicNetworkAccessDisabled
				}
				existing.Properties.PublicNetworkAccess = &publicNetworkAccess
			}

			if metadata.ResourceData.HasChange("sku") {
				existing.Properties.Sku = &model.Sku
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = &model.Tags
			}

			if err := client.AccountsCreateThenPoll(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r IotHubDeviceUpdateAccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.AccountsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
