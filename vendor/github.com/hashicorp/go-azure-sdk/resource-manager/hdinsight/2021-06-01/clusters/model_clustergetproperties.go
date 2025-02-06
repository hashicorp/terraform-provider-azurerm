package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterGetProperties struct {
	ClusterDefinition             ClusterDefinition                  `json:"clusterDefinition"`
	ClusterHdpVersion             *string                            `json:"clusterHdpVersion,omitempty"`
	ClusterId                     *string                            `json:"clusterId,omitempty"`
	ClusterState                  *string                            `json:"clusterState,omitempty"`
	ClusterVersion                *string                            `json:"clusterVersion,omitempty"`
	ComputeIsolationProperties    *ComputeIsolationProperties        `json:"computeIsolationProperties,omitempty"`
	ComputeProfile                *ComputeProfile                    `json:"computeProfile,omitempty"`
	ConnectivityEndpoints         *[]ConnectivityEndpoint            `json:"connectivityEndpoints,omitempty"`
	CreatedDate                   *string                            `json:"createdDate,omitempty"`
	DiskEncryptionProperties      *DiskEncryptionProperties          `json:"diskEncryptionProperties,omitempty"`
	EncryptionInTransitProperties *EncryptionInTransitProperties     `json:"encryptionInTransitProperties,omitempty"`
	Errors                        *[]Errors                          `json:"errors,omitempty"`
	ExcludedServicesConfig        *ExcludedServicesConfig            `json:"excludedServicesConfig,omitempty"`
	KafkaRestProperties           *KafkaRestProperties               `json:"kafkaRestProperties,omitempty"`
	MinSupportedTlsVersion        *string                            `json:"minSupportedTlsVersion,omitempty"`
	NetworkProperties             *NetworkProperties                 `json:"networkProperties,omitempty"`
	OsType                        *OSType                            `json:"osType,omitempty"`
	PrivateEndpointConnections    *[]PrivateEndpointConnection       `json:"privateEndpointConnections,omitempty"`
	PrivateLinkConfigurations     *[]PrivateLinkConfiguration        `json:"privateLinkConfigurations,omitempty"`
	ProvisioningState             *HDInsightClusterProvisioningState `json:"provisioningState,omitempty"`
	QuotaInfo                     *QuotaInfo                         `json:"quotaInfo,omitempty"`
	SecurityProfile               *SecurityProfile                   `json:"securityProfile,omitempty"`
	StorageProfile                *StorageProfile                    `json:"storageProfile,omitempty"`
	Tier                          *Tier                              `json:"tier,omitempty"`
}
