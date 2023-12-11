// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/blobcontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type storageContainersDataSource struct{}

var _ sdk.DataSource = storageContainersDataSource{}

type storageContainersDataSourceModel struct {
	StorageAccountId string           `tfschema:"storage_account_id"`
	NamePrefix       string           `tfschema:"name_prefix"`
	Containers       []containerModel `tfschema:"containers"`
}

type containerModel struct {
	Name              string `tfschema:"name"`
	DataPlaneId       string `tfschema:"data_plane_id"`
	ResourceManagerId string `tfschema:"resource_manager_id"`
}

func (r storageContainersDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},
		"name_prefix": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r storageContainersDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"containers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"data_plane_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"resource_manager_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r storageContainersDataSource) ResourceType() string {
	return "azurerm_storage_containers"
}

func (r storageContainersDataSource) ModelObject() interface{} {
	return &storageContainersDataSourceModel{}
}

func (r storageContainersDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.BlobContainers

			var plan storageContainersDataSourceModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := commonids.ParseStorageAccountID(plan.StorageAccountId)
			if err != nil {
				return err
			}

			resp, err := client.ListCompleteMatchingPredicate(ctx, *id, blobcontainers.DefaultListOperationOptions(), blobcontainers.ListContainerItemOperationPredicate{})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			plan.Containers = flattenStorageContainersContainers(resp.Items, id.StorageAccountName, metadata.Client.Storage.Environment.StorageEndpointSuffix, plan.NamePrefix)

			if err := metadata.Encode(&plan); err != nil {
				return fmt.Errorf("encoding %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func flattenStorageContainersContainers(l []blobcontainers.ListContainerItem, accountName, endpointSuffix, prefix string) []containerModel {
	var output []containerModel
	for _, item := range l {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		if prefix != "" && !strings.HasPrefix(name, prefix) {
			continue
		}

		var mgmtId string
		if item.Id != nil {
			mgmtId = *item.Id
		}

		output = append(output, containerModel{
			Name:              name,
			ResourceManagerId: mgmtId,
			DataPlaneId:       parse.NewStorageContainerDataPlaneId(accountName, endpointSuffix, name).ID(),
		})
	}

	return output
}
