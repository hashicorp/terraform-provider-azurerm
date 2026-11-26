// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type MongoClusterResourceModel struct {
	Name                         string                         `tfschema:"name"`
	ResourceGroupName            string                         `tfschema:"resource_group_name"`
	Location                     string                         `tfschema:"location"`
	AdministratorUserName        string                         `tfschema:"administrator_username"`
	AdministratorPassword        string                         `tfschema:"administrator_password"`
	AuthenticationMethods        []string                       `tfschema:"authentication_methods"`
	CreateMode                   string                         `tfschema:"create_mode"`
	CustomerManagedKey           []CustomerManagedKey           `tfschema:"customer_managed_key"`
	DataApiModeEnabled           bool                           `tfschema:"data_api_mode_enabled"`
	Identity                     []identity.ModelUserAssigned   `tfschema:"identity"`
	Restore                      []Restore                      `tfschema:"restore"`
	ShardCount                   int64                          `tfschema:"shard_count"`
	SourceLocation               string                         `tfschema:"source_location"`
	SourceServerId               string                         `tfschema:"source_server_id"`
	ComputeTier                  string                         `tfschema:"compute_tier"`
	HighAvailabilityMode         string                         `tfschema:"high_availability_mode"`
	CosmosDBNetworkBypassEnabled bool                           `tfschema:"cosmos_db_network_bypass_enabled"`
	PublicNetworkAccess          string                         `tfschema:"public_network_access"`
	PreviewFeatures              []string                       `tfschema:"preview_features"`
	StorageSizeInGb              int64                          `tfschema:"storage_size_in_gb"`
	StorageType                  string                         `tfschema:"storage_type"`
	ConnectionStrings            []MongoClusterConnectionString `tfschema:"connection_strings"`
	Tags                         map[string]string              `tfschema:"tags"`
	Version                      string                         `tfschema:"version"`
}

func (r Resource) ModelObject() any {
	return &MongoClusterResourceModel{}
}
