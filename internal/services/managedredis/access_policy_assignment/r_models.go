// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package access_policy_assignment

type ResourceModel struct {
	ManagedRedisID string `tfschema:"managed_redis_id"`
	ObjectID       string `tfschema:"object_id"`
}

const defaultDatabaseName = "default"

func (r Resource) ModelObject() interface{} {
	return &ResourceModel{}
}
