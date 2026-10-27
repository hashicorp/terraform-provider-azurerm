// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache_access_policy

type RedisCacheAccessPolicyResourceModel struct {
	Name         string `tfschema:"name"`
	RedisCacheID string `tfschema:"redis_cache_id"`
	Permissions  string `tfschema:"permissions"`
}

func (r Resource) ModelObject() any {
	return &RedisCacheAccessPolicyResourceModel{}
}
