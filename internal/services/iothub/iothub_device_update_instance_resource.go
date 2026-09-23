// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type IotHubDeviceUpdateInstanceResource struct{}

var _ sdk.ResourceWithUpdate = IotHubDeviceUpdateInstanceResource{}

type IotHubDeviceUpdateInstanceModel struct {
	Name                     string                          `tfschema:"name"`
	DeviceUpdateAccountId    string                          `tfschema:"device_update_account_id"`
	DiagnosticStorageAccount []DiagnosticStorageAccountModel `tfschema:"diagnostic_storage_account"`
	DiagnosticEnabled        bool                            `tfschema:"diagnostic_enabled"`
	IotHubId                 string                          `tfschema:"iothub_id"`
	Tags                     map[string]string               `tfschema:"tags"`
}

type DiagnosticStorageAccountModel struct {
	ConnectionString          string `tfschema:"connection_string"`
	ConnectionStringWOVersion int64  `tfschema:"connection_string_wo_version"`
	Id                        string `tfschema:"id"`
}

func (r IotHubDeviceUpdateInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubDeviceUpdateInstanceName,
		},

		"device_update_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: deviceupdates.ValidateAccountID,
		},

		"iothub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubID,
		},

		"diagnostic_storage_account": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"connection_string": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"diagnostic_storage_account.0.connection_string_wo"},
					},

					"connection_string_wo": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						WriteOnly:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"diagnostic_storage_account.0.connection_string"},
						RequiredWith: []string{"diagnostic_storage_account.0.connection_string_wo_version"},
					},

					"connection_string_wo_version": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
						RequiredWith: []string{"diagnostic_storage_account.0.connection_string_wo"},
					},
				},
			},
		},

		"diagnostic_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}

func (r IotHubDeviceUpdateInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r IotHubDeviceUpdateInstanceResource) ResourceType() string {
	return "azurerm_iothub_device_update_instance"
}

func (r IotHubDeviceUpdateInstanceResource) ModelObject() interface{} {
	return &IotHubDeviceUpdateInstanceModel{}
}

func (r IotHubDeviceUpdateInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deviceupdates.ValidateInstanceID
}

func (r IotHubDeviceUpdateInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model IotHubDeviceUpdateInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.IoTHub.DeviceUpdatesClient
			deviceUpdateAccountId, err := deviceupdates.ParseAccountID(model.DeviceUpdateAccountId)
			if err != nil {
				return err
			}

			id := deviceupdates.NewInstanceID(deviceUpdateAccountId.SubscriptionId, deviceUpdateAccountId.ResourceGroupName, deviceUpdateAccountId.AccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.InstancesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			deviceUpdateAccount, err := client.AccountsGet(ctx, *deviceUpdateAccountId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *deviceUpdateAccountId, err)
			}

			if deviceUpdateAccount.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *deviceUpdateAccountId)
			}

			expandedStorageProps, err := expandDiagnosticStorageAccount(metadata, model.DiagnosticStorageAccount)
			if err != nil {
				return fmt.Errorf("expanding `diagnostic_storage_account`: %w", err)
			}

			properties := &deviceupdates.Instance{
				Location: location.Normalize(deviceUpdateAccount.Model.Location),
				Properties: deviceupdates.InstanceProperties{
					AccountName:                 &deviceUpdateAccountId.AccountName,
					DiagnosticStorageProperties: expandedStorageProps,
					EnableDiagnostics:           &model.DiagnosticEnabled,
					IotHubs: &[]deviceupdates.IotHubSettings{
						{
							ResourceId: model.IotHubId,
						},
					},
				},
				Tags: &model.Tags,
			}

			if err := client.InstancesCreateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r IotHubDeviceUpdateInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.InstancesGet(ctx, *id)
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

			state := IotHubDeviceUpdateInstanceModel{
				Name:                  *model.Name,
				DeviceUpdateAccountId: deviceupdates.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			properties := model.Properties

			if iotHubs := properties.IotHubs; iotHubs != nil && len(*iotHubs) > 0 {
				state.IotHubId = (*iotHubs)[0].ResourceId
			}

			state.DiagnosticStorageAccount = flattenDiagnosticStorageAccount(properties.DiagnosticStorageProperties, metadata)

			state.DiagnosticEnabled = pointer.From(properties.EnableDiagnostics)

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r IotHubDeviceUpdateInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model IotHubDeviceUpdateInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.InstancesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			existing := resp.Model
			if existing == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			// connectionString is not returned by API, so always expands DiagnosticStorageAccount
			expandedStorageProps, err := expandDiagnosticStorageAccount(metadata, model.DiagnosticStorageAccount)
			if err != nil {
				return fmt.Errorf("expanding `diagnostic_storage_account`: %w", err)
			}
			existing.Properties.DiagnosticStorageProperties = expandedStorageProps

			if metadata.ResourceData.HasChange("diagnostic_enabled") {
				existing.Properties.EnableDiagnostics = &model.DiagnosticEnabled
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = &model.Tags
			}

			if err := client.InstancesCreateThenPoll(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r IotHubDeviceUpdateInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.InstancesDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandDiagnosticStorageAccount(rmd sdk.ResourceMetaData, inputList []DiagnosticStorageAccountModel) (*deviceupdates.DiagnosticStorageProperties, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	woConnectionString, err := pluginsdk.GetWriteOnly(rmd.ResourceData, "diagnostic_storage_account.0.connection_string_wo", cty.String)
	if err != nil {
		return nil, err
	}

	input := inputList[0]
	output := deviceupdates.DiagnosticStorageProperties{
		AuthenticationType: deviceupdates.AuthenticationTypeKeyBased,
		ConnectionString:   &input.ConnectionString,
		ResourceId:         input.Id,
	}

	if !woConnectionString.IsNull() {
		output.ConnectionString = pointer.To(woConnectionString.AsString())
	}

	return &output, nil
}

func flattenDiagnosticStorageAccount(input *deviceupdates.DiagnosticStorageProperties, metadata sdk.ResourceMetaData) []DiagnosticStorageAccountModel {
	var outputList []DiagnosticStorageAccountModel
	if input == nil {
		return []DiagnosticStorageAccountModel{}
	}

	output := DiagnosticStorageAccountModel{
		Id:                        input.ResourceId,
		ConnectionStringWOVersion: int64(metadata.ResourceData.Get("diagnostic_storage_account.0.connection_string_wo_version").(int)),
	}

	// connectionString is not returned by API
	if connectionString, ok := metadata.ResourceData.GetOk("diagnostic_storage_account.0.connection_string"); ok {
		output.ConnectionString = connectionString.(string)
	}

	return append(outputList, output)
}
