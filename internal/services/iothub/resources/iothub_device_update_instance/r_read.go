// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_device_update_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
