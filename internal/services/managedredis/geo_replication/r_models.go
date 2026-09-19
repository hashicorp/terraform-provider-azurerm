// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_replication

type ResourceModel struct {
	ManagedRedisId        string   `tfschema:"managed_redis_id"`
	LinkedManagedRedisIds []string `tfschema:"linked_managed_redis_ids"`
}

const defaultDatabaseName = "default"

func (r Resource) ModelObject() interface{} {
	return &ResourceModel{}
}
