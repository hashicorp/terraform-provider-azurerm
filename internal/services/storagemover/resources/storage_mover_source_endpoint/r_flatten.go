// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_source_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *endpoints.EndpointId, model *endpoints.Endpoint) error {
	state := StorageMoverSourceEndpointModel{
		Name:           id.EndpointName,
		StorageMoverId: storagemovers.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
	}

	if model != nil {
		if v, ok := model.Properties.(endpoints.NfsMountEndpointProperties); ok {
			state.Export = v.Export
			state.Host = v.Host

			if v := v.NfsVersion; v != nil {
				state.NfsVersion = *v
			}

			state.Description = pointer.From(v.Description)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
