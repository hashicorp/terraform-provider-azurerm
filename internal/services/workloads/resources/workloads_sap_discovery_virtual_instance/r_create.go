// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_discovery_virtual_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPDiscoveryVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPVirtualInstances
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := sapvirtualinstances.NewSapVirtualInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			identity, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := &sapvirtualinstances.SAPVirtualInstance{
				Identity: identity,
				Location: location.Normalize(model.Location),
				Properties: &sapvirtualinstances.SAPVirtualInstanceProperties{
					Environment:                       sapvirtualinstances.SAPEnvironmentType(model.Environment),
					SapProduct:                        sapvirtualinstances.SAPProductType(model.SapProduct),
					ManagedResourcesNetworkAccessType: pointer.ToEnum[sapvirtualinstances.ManagedResourcesNetworkAccessType](model.ManagedResourcesNetworkAccessType),
				},
				Tags: &model.Tags,
			}

			discoveryConfiguration := &sapvirtualinstances.DiscoveryConfiguration{
				CentralServerVMId: pointer.To(model.CentralServerVmId),
			}

			if v := model.ManagedStorageAccountName; v != "" {
				discoveryConfiguration.ManagedRgStorageAccountName = pointer.To(v)
			}

			parameters.Properties.Configuration = discoveryConfiguration

			if v := model.ManagedResourceGroupName; v != "" {
				parameters.Properties.ManagedResourceGroupConfiguration = &sapvirtualinstances.ManagedRGConfiguration{
					Name: pointer.To(v),
				}
			}

			if err := client.CreateCallbackThenPoll(ctx, id, *parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err = pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return fmt.Errorf("setting resource identity data: %+v", err)
			}
			return nil
		},
	}
}
