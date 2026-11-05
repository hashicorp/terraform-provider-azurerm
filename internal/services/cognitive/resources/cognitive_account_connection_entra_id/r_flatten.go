// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_entra_id

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *accountconnectionresource.ConnectionId, model *accountconnectionresource.ConnectionPropertiesV2BasicResource, priorMetadata map[string]string) error {
	state := CognitiveAccountConnectionEntraIDModel{
		CognitiveAccountId: accountconnectionresource.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
		Name:               id.ConnectionName,
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	if model != nil && model.Properties != nil {
		base := model.Properties.ConnectionPropertiesV2()
		state.Category = pointer.FromEnum(base.Category)
		state.Target = pointer.From(base.Target)
		state.Metadata = helpers.FlattenAccountConnectionMetadata(priorMetadata, base.Metadata)
	}

	return metadata.Encode(&state)
}
