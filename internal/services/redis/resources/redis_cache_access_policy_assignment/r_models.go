// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache_access_policy_assignment

type RedisCacheAccessPolicyAssignmentResourceModel struct {
	Name             string `tfschema:"name"`
	RedisCacheID     string `tfschema:"redis_cache_id"`
	AccessPolicyName string `tfschema:"access_policy_name"`
	ObjectID         string `tfschema:"object_id"`
	ObjectIDAlias    string `tfschema:"object_id_alias"`
}

func (r Resource) ModelObject() any {
	return &RedisCacheAccessPolicyAssignmentResourceModel{}
}
