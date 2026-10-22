// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package federated_identity_credential

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *federatedidentitycredentials.FederatedIdentityCredentialId, model *federatedidentitycredentials.FederatedIdentityCredential) error {
	schema := FederatedIdentityCredentialResourceSchema{
		Name: id.FederatedIdentityCredentialName,
	}

	parentId := commonids.NewUserAssignedIdentityID(id.SubscriptionId, id.ResourceGroupName, id.UserAssignedIdentityName)
	schema.UserAssignedIdentityId = parentId.ID()

	if model != nil {
		r.mapFederatedIdentityCredentialToFederatedIdentityCredentialResourceSchema(*model, &schema)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&schema)
}
