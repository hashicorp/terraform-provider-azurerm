// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

const (
	DatabaseModuleRedisBloom      = "RedisBloom"
	DatabaseModuleRediSearch      = "RediSearch"
	DatabaseModuleRedisJSON       = "RedisJSON"
	DatabaseModuleRedisTimeSeries = "RedisTimeSeries"
)

func DatabaseModulesSupportingGeoReplication() []string {
	return []string{
		DatabaseModuleRediSearch,
		DatabaseModuleRedisJSON,
	}
}
