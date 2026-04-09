package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterCreateProperties struct {
	ClusterDefinition             *ClusterDefinition             `json:"clusterDefinition,omitempty"`
	ClusterVersion                *string                        `json:"clusterVersion,omitempty"`
	ComputeIsolationProperties    *ComputeIsolationProperties    `json:"computeIsolationProperties,omitempty"`
	ComputeProfile                *ComputeProfile                `json:"computeProfile,omitempty"`
	DiskEncryptionProperties      *DiskEncryptionProperties      `json:"diskEncryptionProperties,omitempty"`
	EncryptionInTransitProperties *EncryptionInTransitProperties `json:"encryptionInTransitProperties,omitempty"`
	KafkaRestProperties           *KafkaRestProperties           `json:"kafkaRestProperties,omitempty"`
	MinSupportedTlsVersion        *string                        `json:"minSupportedTlsVersion,omitempty"`
	NetworkProperties             *NetworkProperties             `json:"networkProperties,omitempty"`
	OsType                        *OSType                        `json:"osType,omitempty"`
	PrivateLinkConfigurations     *[]PrivateLinkConfiguration    `json:"privateLinkConfigurations,omitempty"`
	SecurityProfile               *SecurityProfile               `json:"securityProfile,omitempty"`
	StorageProfile                *StorageProfile                `json:"storageProfile,omitempty"`
	Tier                          *Tier                          `json:"tier,omitempty"`
}
