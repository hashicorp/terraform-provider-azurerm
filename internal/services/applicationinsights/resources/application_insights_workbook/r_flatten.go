// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *workbooksapis.WorkbookId, model *workbooksapis.Workbook) error {
	state := ApplicationInsightsWorkbookModel{
		Name:              id.WorkbookName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if properties := model.Properties; properties != nil {
			state.Category = properties.Category

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			state.DisplayName = properties.DisplayName

			state.DataJson = properties.SerializedData

			if properties.SourceId != nil {
				state.SourceId = *properties.SourceId
			}

			if properties.StorageUri != nil {
				state.StorageContainerId = *properties.StorageUri
			}
		}

		if model.Tags != nil {
			// The backend returns a tags with key `hidden-title` by default. Since it has the same value with `display_name` and will cause inconsistency with user's configuration, remove it as a workaround.
			delete(*model.Tags, "hidden-title")
			state.Tags = *model.Tags
		}
	}
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
