// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_catalog

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbitalplanetarycomputer/2026-04-15/geocatalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *geocatalogs.GeoCatalogId, model *geocatalogs.GeoCatalog) error {
	state := GeoCatalogModel{
		Name:              id.GeoCatalogName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = flattenedIdentity

		state.Tags = pointer.From(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
