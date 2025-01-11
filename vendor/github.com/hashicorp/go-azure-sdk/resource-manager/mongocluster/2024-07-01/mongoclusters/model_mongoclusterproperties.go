package mongoclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoClusterProperties struct {
	Administrator              *AdministratorProperties       `json:"administrator,omitempty"`
	Backup                     *BackupProperties              `json:"backup,omitempty"`
	ClusterStatus              *MongoClusterStatus            `json:"clusterStatus,omitempty"`
	Compute                    *ComputeProperties             `json:"compute,omitempty"`
	ConnectionString           *string                        `json:"connectionString,omitempty"`
	CreateMode                 *CreateMode                    `json:"createMode,omitempty"`
	HighAvailability           *HighAvailabilityProperties    `json:"highAvailability,omitempty"`
	InfrastructureVersion      *string                        `json:"infrastructureVersion,omitempty"`
	PreviewFeatures            *[]PreviewFeature              `json:"previewFeatures,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection   `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState             `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess           `json:"publicNetworkAccess,omitempty"`
	Replica                    *ReplicationProperties         `json:"replica,omitempty"`
	ReplicaParameters          *MongoClusterReplicaParameters `json:"replicaParameters,omitempty"`
	RestoreParameters          *MongoClusterRestoreParameters `json:"restoreParameters,omitempty"`
	ServerVersion              *string                        `json:"serverVersion,omitempty"`
	Sharding                   *ShardingProperties            `json:"sharding,omitempty"`
	Storage                    *StorageProperties             `json:"storage,omitempty"`
}
