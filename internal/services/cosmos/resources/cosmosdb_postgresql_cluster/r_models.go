// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_cluster

type CosmosDbPostgreSQLClusterModel struct {
	Name                             string              `tfschema:"name"`
	ResourceGroupName                string              `tfschema:"resource_group_name"`
	Location                         string              `tfschema:"location"`
	AdministratorLoginPassword       string              `tfschema:"administrator_login_password"`
	CitusVersion                     string              `tfschema:"citus_version"`
	CoordinatorPublicIPAccessEnabled bool                `tfschema:"coordinator_public_ip_access_enabled"`
	CoordinatorServerEdition         string              `tfschema:"coordinator_server_edition"`
	CoordinatorStorageQuotaInMb      int64               `tfschema:"coordinator_storage_quota_in_mb"`
	CoordinatorVCoreCount            int64               `tfschema:"coordinator_vcore_count"`
	HaEnabled                        bool                `tfschema:"ha_enabled"`
	ShardsOnCoordinatorEnabled       bool                `tfschema:"shards_on_coordinator_enabled"`
	SourceLocation                   string              `tfschema:"source_location"`
	SourceResourceId                 string              `tfschema:"source_resource_id"`
	MaintenanceWindow                []MaintenanceWindow `tfschema:"maintenance_window"`
	ServerNames                      []ServerNameItem    `tfschema:"servers"`
	NodeCount                        int64               `tfschema:"node_count"`
	NodePublicIPAccessEnabled        bool                `tfschema:"node_public_ip_access_enabled"`
	NodeServerEdition                string              `tfschema:"node_server_edition"`
	NodeStorageQuotaInMb             int64               `tfschema:"node_storage_quota_in_mb"`
	NodeVCores                       int64               `tfschema:"node_vcores"`
	PointInTimeInUTC                 string              `tfschema:"point_in_time_in_utc"`
	PreferredPrimaryZone             string              `tfschema:"preferred_primary_zone"`
	SqlVersion                       string              `tfschema:"sql_version"`
	Tags                             map[string]string   `tfschema:"tags"`
	EarliestRestoreTime              string              `tfschema:"earliest_restore_time"`
}

func (r Resource) ModelObject() any {
	return &CosmosDbPostgreSQLClusterModel{}
}
