// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type AIFoundryModel struct {
	Name                        string                                     `tfschema:"name"`
	Location                    string                                     `tfschema:"location"`
	ResourceGroupName           string                                     `tfschema:"resource_group_name"`
	ApplicationInsightsId       string                                     `tfschema:"application_insights_id"`
	StorageAccountId            string                                     `tfschema:"storage_account_id"`
	KeyVaultId                  string                                     `tfschema:"key_vault_id"`
	ContainerRegistryId         string                                     `tfschema:"container_registry_id"`
	Encryption                  []Encryption                               `tfschema:"encryption"`
	ManagedNetwork              []ManagedNetwork                           `tfschema:"managed_network"`
	PublicNetworkAccess         string                                     `tfschema:"public_network_access"`
	Identity                    []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PrimaryUserAssignedIdentity string                                     `tfschema:"primary_user_assigned_identity"`
	HighBusinessImpactEnabled   bool                                       `tfschema:"high_business_impact_enabled"`
	Description                 string                                     `tfschema:"description"`
	FriendlyName                string                                     `tfschema:"friendly_name"`
	DiscoveryUrl                string                                     `tfschema:"discovery_url"`
	WorkspaceId                 string                                     `tfschema:"workspace_id"`
	Tags                        map[string]any                             `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &AIFoundryModel{}
}
