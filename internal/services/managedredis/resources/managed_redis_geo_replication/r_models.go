// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis_geo_replication

type ManagedRedisGeoReplicationResourceModel struct {
	ManagedRedisId        string   `tfschema:"managed_redis_id"`
	LinkedManagedRedisIds []string `tfschema:"linked_managed_redis_ids"`
}

func (r Resource) ModelObject() interface{} {
	return &ManagedRedisGeoReplicationResourceModel{}
}
