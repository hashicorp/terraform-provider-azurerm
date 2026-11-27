// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_bookshelf

type DiscoveryBookshelfModel struct {
	Name                       string                                 `tfschema:"name"`
	ResourceGroupName          string                                 `tfschema:"resource_group_name"`
	Location                   string                                 `tfschema:"location"`
	CustomerManagedKey         []DiscoveryBookshelfCustomerManagedKey `tfschema:"customer_managed_key"`
	LogAnalyticsClusterId      string                                 `tfschema:"log_analytics_cluster_id"`
	PrivateEndpointSubnetId    string                                 `tfschema:"private_endpoint_subnet_id"`
	PublicNetworkAccessEnabled bool                                   `tfschema:"public_network_access_enabled"`
	SearchSubnetId             string                                 `tfschema:"search_subnet_id"`
	WorkloadIdentityIds        []string                               `tfschema:"workload_identity_ids"`
	Tags                       map[string]string                      `tfschema:"tags"`

	ManagedOnBehalfOfBrokerResourceIds []string `tfschema:"managed_on_behalf_of_broker_resource_ids"`
	ManagedResourceGroupName           string   `tfschema:"managed_resource_group_name"`
	Uri                                string   `tfschema:"uri"`
}

type DiscoveryBookshelfCustomerManagedKey struct {
	KeyVaultKeyId          string `tfschema:"key_vault_key_id"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
}

func (r Resource) ModelObject() any {
	return &DiscoveryBookshelfModel{}
}
