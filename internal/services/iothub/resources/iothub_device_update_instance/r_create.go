// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_device_update_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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
