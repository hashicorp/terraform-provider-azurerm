package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisCommonPropertiesRedisConfiguration struct {
	AofBackupEnabled                   *string `json:"aof-backup-enabled,omitempty"`
	AofStorageConnectionString0        *string `json:"aof-storage-connection-string-0,omitempty"`
	AofStorageConnectionString1        *string `json:"aof-storage-connection-string-1,omitempty"`
	Authnotrequired                    *string `json:"authnotrequired,omitempty"`
	Maxclients                         *string `json:"maxclients,omitempty"`
	MaxfragmentationmemoryReserved     *string `json:"maxfragmentationmemory-reserved,omitempty"`
	MaxmemoryDelta                     *string `json:"maxmemory-delta,omitempty"`
	MaxmemoryPolicy                    *string `json:"maxmemory-policy,omitempty"`
	MaxmemoryReserved                  *string `json:"maxmemory-reserved,omitempty"`
	NotifyKeyspaceEvents               *string `json:"notify-keyspace-events,omitempty"`
	PreferredDataArchiveAuthMethod     *string `json:"preferred-data-archive-auth-method,omitempty"`
	PreferredDataPersistenceAuthMethod *string `json:"preferred-data-persistence-auth-method,omitempty"`
	RdbBackupEnabled                   *string `json:"rdb-backup-enabled,omitempty"`
	RdbBackupFrequency                 *string `json:"rdb-backup-frequency,omitempty"`
	RdbBackupMaxSnapshotCount          *string `json:"rdb-backup-max-snapshot-count,omitempty"`
	RdbStorageConnectionString         *string `json:"rdb-storage-connection-string,omitempty"`
	StorageSubscriptionId              *string `json:"storage-subscription-id,omitempty"`
	ZonalConfiguration                 *string `json:"zonal-configuration,omitempty"`
}
