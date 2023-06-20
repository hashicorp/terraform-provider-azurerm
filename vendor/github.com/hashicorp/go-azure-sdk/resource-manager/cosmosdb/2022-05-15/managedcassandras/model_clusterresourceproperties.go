package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterResourceProperties struct {
	AuthenticationMethod          *AuthenticationMethod              `json:"authenticationMethod,omitempty"`
	CassandraAuditLoggingEnabled  *bool                              `json:"cassandraAuditLoggingEnabled,omitempty"`
	CassandraVersion              *string                            `json:"cassandraVersion,omitempty"`
	ClientCertificates            *[]Certificate                     `json:"clientCertificates,omitempty"`
	ClusterNameOverride           *string                            `json:"clusterNameOverride,omitempty"`
	Deallocated                   *bool                              `json:"deallocated,omitempty"`
	DelegatedManagementSubnetId   *string                            `json:"delegatedManagementSubnetId,omitempty"`
	ExternalGossipCertificates    *[]Certificate                     `json:"externalGossipCertificates,omitempty"`
	ExternalSeedNodes             *[]SeedNode                        `json:"externalSeedNodes,omitempty"`
	GossipCertificates            *[]Certificate                     `json:"gossipCertificates,omitempty"`
	HoursBetweenBackups           *int64                             `json:"hoursBetweenBackups,omitempty"`
	InitialCassandraAdminPassword *string                            `json:"initialCassandraAdminPassword,omitempty"`
	PrometheusEndpoint            *SeedNode                          `json:"prometheusEndpoint,omitempty"`
	ProvisioningState             *ManagedCassandraProvisioningState `json:"provisioningState,omitempty"`
	RepairEnabled                 *bool                              `json:"repairEnabled,omitempty"`
	RestoreFromBackupId           *string                            `json:"restoreFromBackupId,omitempty"`
	SeedNodes                     *[]SeedNode                        `json:"seedNodes,omitempty"`
}
