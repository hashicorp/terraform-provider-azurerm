// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			var model SystemCenterVirtualMachineManagerServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := vmmservers.NewVMmServerID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := &vmmservers.VMmServer{
				Location: location.Normalize(model.Location),
				ExtendedLocation: vmmservers.ExtendedLocation{
					Type: pointer.To("customLocation"),
					Name: pointer.To(model.CustomLocationId),
				},
				Properties: &vmmservers.VMmServerProperties{
					Credentials: &vmmservers.VMmCredential{
						Username: pointer.To(model.Username),
						Password: pointer.To(model.Password),
					},
					Fqdn: model.Fqdn,
				},
				Tags: pointer.To(model.Tags),
			}

			if v := model.Port; v != 0 {
				parameters.Properties.Port = pointer.To(v)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			// After System Center Virtual Machine Manager Server is created, it needs some time to sync the Inventory Items. And service team confirmed that the sync would definitely be completed within 10 minutes. In case, so we need to set a timeout of 120 minutes and check the inventory quantity continuously every minute for 10 times. If the quantity doesn't change, then we consider the sync to be complete.
			stateConf := &pluginsdk.StateChangeConf{
				Delay:        5 * time.Second,
				Pending:      []string{"SyncNotCompleted"},
				Target:       []string{"SyncCompleted"},
				Refresh:      systemCenterVirtualMachineManagerServerStateRefreshFunc(ctx, metadata, id),
				PollInterval: 1 * time.Minute,
				Timeout:      120 * time.Minute,
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}
