// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

const (
	DatabaseModuleRedisBloom      = "RedisBloom"
	DatabaseModuleRediSearch      = "RediSearch"
	DatabaseModuleRedisJSON       = "RedisJSON"
	DatabaseModuleRedisTimeSeries = "RedisTimeSeries"
)

func PossibleValuesForDatabaseModule() []string {
	return []string{
		DatabaseModuleRedisBloom,
		DatabaseModuleRediSearch,
		DatabaseModuleRedisJSON,
		DatabaseModuleRedisTimeSeries,
	}
}

func DatabaseModulesSupportingGeoReplication() []string {
	return []string{
		DatabaseModuleRediSearch,
		DatabaseModuleRedisJSON,
	}
}
