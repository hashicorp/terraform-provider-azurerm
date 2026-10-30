// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_entra_id

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *projectconnectionresource.ProjectConnectionId, model *projectconnectionresource.ConnectionPropertiesV2BasicResource, priorMetadata map[string]string) error {
	state := CognitiveAccountProjectConnectionEntraIDModel{
		CognitiveAccountProjectId: projectconnectionresource.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ProjectName).ID(),
		Name:                      id.ConnectionName,
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	if model != nil && model.Properties != nil {
		base := model.Properties.ConnectionPropertiesV2()
		state.AuthType = string(base.AuthType)
		state.Category = pointer.FromEnum(base.Category)
		state.Target = pointer.From(base.Target)

		state.Metadata = helpers.FlattenAccountProjectConnectionMetadata(priorMetadata, base.Metadata)
	}

	return metadata.Encode(&state)
}
