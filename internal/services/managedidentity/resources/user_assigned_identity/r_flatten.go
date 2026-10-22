// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package user_assigned_identity

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/identities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *commonids.UserAssignedIdentityId, model *identities.Identity) error {
	state := UserAssignedIdentityResourceSchema{
		Name:              id.UserAssignedIdentityName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = tags.Flatten(model.Tags)

		if model.Properties != nil {
			state.ClientId = pointer.From(model.Properties.ClientId)
			state.PrincipalId = pointer.From(model.Properties.PrincipalId)
			state.TenantId = pointer.From(model.Properties.TenantId)

			if isolationScope := pointer.FromEnum(model.Properties.IsolationScope); isolationScope != string(identities.IsolationScopeNone) {
				state.IsolationScope = isolationScope
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
