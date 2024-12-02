package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbMigrationSettings struct {
	BoostRUs    *int64                             `json:"boostRUs,omitempty"`
	Databases   map[string]MongoDbDatabaseSettings `json:"databases"`
	Replication *MongoDbReplication                `json:"replication,omitempty"`
	Source      MongoDbConnectionInfo              `json:"source"`
	Target      MongoDbConnectionInfo              `json:"target"`
	Throttling  *MongoDbThrottlingSettings         `json:"throttling,omitempty"`
}
