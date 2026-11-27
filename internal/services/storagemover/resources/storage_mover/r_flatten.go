// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *storagemovers.StorageMoverId, model *storagemovers.StorageMover) error {
	state := StorageMoverModel{
		Name:              id.StorageMoverName,
		ResourceGroupName: id.ResourceGroupName,
		Location:          location.Normalize(model.Location),
	}

	description := ""
	if properties := model.Properties; properties != nil {
		if properties.Description != nil {
			description = *properties.Description
		}
	}
	state.Description = description

	if model.Tags != nil {
		state.Tags = *model.Tags
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
