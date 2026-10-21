// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databox_edge_device

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.DataboxEdge.DeviceClient

			var metaModel EdgeDeviceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := devices.NewDataBoxEdgeDeviceID(subscriptionId, metaModel.ResourceGroupName, metaModel.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_databox_edge_device", id.ID())
				}
			}

			dataBoxEdgeDevice := devices.DataBoxEdgeDevice{
				Location: location.Normalize(metaModel.Location),
				Sku:      expandDeviceSku(metaModel.SkuName),
				Tags:     &metaModel.Tags,
			}

			if _, err := client.CreateOrUpdate(ctx, id, dataBoxEdgeDevice); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func expandDeviceSku(input string) *devices.Sku {
	if len(input) == 0 {
		return nil
	}

	v, err := parse.DataboxEdgeDeviceSkuName(input)
	if err != nil {
		return nil
	}

	return &devices.Sku{
		Name: pointer.ToEnum[devices.SkuName](v.Name),
		Tier: pointer.ToEnum[devices.SkuTier](v.Tier),
	}
}
