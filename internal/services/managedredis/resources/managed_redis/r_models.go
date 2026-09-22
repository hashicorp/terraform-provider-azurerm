// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ManagedRedisResourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

	Location string `tfschema:"location"`

	SkuName string `tfschema:"sku_name"`

	CustomerManagedKey      []CustomerManagedKeyModel                  `tfschema:"customer_managed_key"`
	DefaultDatabase         []DefaultDatabaseModel                     `tfschema:"default_database"`
	HighAvailabilityEnabled bool                                       `tfschema:"high_availability_enabled"`
	Identity                []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccess     string                                     `tfschema:"public_network_access"`
	Tags                    map[string]string                          `tfschema:"tags"`

	Hostname string `tfschema:"hostname"`
}

type CustomerManagedKeyModel struct {
	KeyVaultKeyId          string `tfschema:"key_vault_key_id"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
}

type DefaultDatabaseModel struct {
	AccessKeysAuthenticationEnabled          bool          `tfschema:"access_keys_authentication_enabled"`
	ClientProtocol                           string        `tfschema:"client_protocol"`
	ClusteringPolicy                         string        `tfschema:"clustering_policy"`
	EvictionPolicy                           string        `tfschema:"eviction_policy"`
	GeoReplicationGroupName                  string        `tfschema:"geo_replication_group_name"`
	Module                                   []ModuleModel `tfschema:"module"`
	PersistenceAppendOnlyFileBackupFrequency string        `tfschema:"persistence_append_only_file_backup_frequency"`
	PersistenceRedisDatabaseBackupFrequency  string        `tfschema:"persistence_redis_database_backup_frequency"`

	ID                 string `tfschema:"id"`
	Port               int64  `tfschema:"port"`
	PrimaryAccessKey   string `tfschema:"primary_access_key"`
	SecondaryAccessKey string `tfschema:"secondary_access_key"`
}

type ModuleModel struct {
	Name    string `tfschema:"name"`
	Args    string `tfschema:"args"`
	Version string `tfschema:"version"`
}

func (r Resource) ModelObject() interface{} {
	return &ManagedRedisResourceModel{}
}
