// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis_access_policy_assignment

type ManagedRedisAccessPolicyAssignmentResourceModel struct {
	ManagedRedisID string `tfschema:"managed_redis_id"`
	ObjectID       string `tfschema:"object_id"`
}

func (r Resource) ModelObject() any {
	return &ManagedRedisAccessPolicyAssignmentResourceModel{}
}
