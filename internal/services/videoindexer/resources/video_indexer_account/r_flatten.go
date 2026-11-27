// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package video_indexer_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2025-04-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenStorage(input *accounts.StorageServicesForPutRequest) ([]StorageModel, error) {
	if input == nil {
		return []StorageModel{}, nil
	}

	storage := StorageModel{}
	if v := pointer.From(input.ResourceId); v != "" {
		id, err := commonids.ParseStorageAccountID(v)
		if err != nil {
			return []StorageModel{}, err
		}
		storage.StorageAccountId = id.ID()
	}

	if v := pointer.From(input.UserAssignedIdentity); v != "" {
		id, err := commonids.ParseUserAssignedIdentityID(v)
		if err != nil {
			return []StorageModel{}, err
		}
		storage.UserAssignedIdentityId = id.ID()
	}

	return []StorageModel{storage}, nil
}

func (Resource) flatten(metadata sdk.ResourceMetaData, id *accounts.AccountId, model *accounts.Account) error {
	state := AccountModel{
		Name:          id.AccountName,
		ResourceGroup: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = flattenedIdentity

		if props := model.Properties; props != nil {
			state.Storage, err = flattenStorage(props.StorageServices)
			if err != nil {
				return fmt.Errorf("flattening `storage`: %+v", err)
			}

			state.PublicNetworkAccess = pointer.FromEnum(props.PublicNetworkAccess)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
