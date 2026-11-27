// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	elasticsanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := elasticsanClient.From(metadata.Client).Volumes

			var config ElasticSANVolumeResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := account.From(metadata.Client).SubscriptionId

			volumeGroupId, err := volumes.ParseVolumeGroupID(config.VolumeGroupId)
			if err != nil {
				return err
			}

			id := volumes.NewVolumeID(subscriptionId, volumeGroupId.ResourceGroupName, volumeGroupId.ElasticSanName, volumeGroupId.VolumeGroupName, config.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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

			payload := volumes.Volume{
				Properties: volumes.VolumeProperties{
					CreationData: ExpandElasticSANVolumeCreateSource(config.CreateSource),
					SizeGiB:      config.SizeInGiB,
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
