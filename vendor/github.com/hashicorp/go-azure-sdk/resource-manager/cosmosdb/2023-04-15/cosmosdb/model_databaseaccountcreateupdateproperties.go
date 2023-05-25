package cosmosdb

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseAccountCreateUpdateProperties struct {
	AnalyticalStorageConfiguration     *AnalyticalStorageConfiguration `json:"analyticalStorageConfiguration,omitempty"`
	ApiProperties                      *ApiProperties                  `json:"apiProperties,omitempty"`
	BackupPolicy                       BackupPolicy                    `json:"backupPolicy"`
	Capabilities                       *[]Capability                   `json:"capabilities,omitempty"`
	Capacity                           *Capacity                       `json:"capacity,omitempty"`
	ConnectorOffer                     *ConnectorOffer                 `json:"connectorOffer,omitempty"`
	ConsistencyPolicy                  *ConsistencyPolicy              `json:"consistencyPolicy,omitempty"`
	Cors                               *[]CorsPolicy                   `json:"cors,omitempty"`
	CreateMode                         *CreateMode                     `json:"createMode,omitempty"`
	DatabaseAccountOfferType           DatabaseAccountOfferType        `json:"databaseAccountOfferType"`
	DefaultIdentity                    *string                         `json:"defaultIdentity,omitempty"`
	DisableKeyBasedMetadataWriteAccess *bool                           `json:"disableKeyBasedMetadataWriteAccess,omitempty"`
	DisableLocalAuth                   *bool                           `json:"disableLocalAuth,omitempty"`
	EnableAnalyticalStorage            *bool                           `json:"enableAnalyticalStorage,omitempty"`
	EnableAutomaticFailover            *bool                           `json:"enableAutomaticFailover,omitempty"`
	EnableCassandraConnector           *bool                           `json:"enableCassandraConnector,omitempty"`
	EnableFreeTier                     *bool                           `json:"enableFreeTier,omitempty"`
	EnableMultipleWriteLocations       *bool                           `json:"enableMultipleWriteLocations,omitempty"`
	EnablePartitionMerge               *bool                           `json:"enablePartitionMerge,omitempty"`
	IPRules                            *[]IPAddressOrRange             `json:"ipRules,omitempty"`
	IsVirtualNetworkFilterEnabled      *bool                           `json:"isVirtualNetworkFilterEnabled,omitempty"`
	KeyVaultKeyUri                     *string                         `json:"keyVaultKeyUri,omitempty"`
	KeysMetadata                       *DatabaseAccountKeysMetadata    `json:"keysMetadata,omitempty"`
	Locations                          []Location                      `json:"locations"`
	MinimalTlsVersion                  *MinimalTlsVersion              `json:"minimalTlsVersion,omitempty"`
	NetworkAclBypass                   *NetworkAclBypass               `json:"networkAclBypass,omitempty"`
	NetworkAclBypassResourceIds        *[]string                       `json:"networkAclBypassResourceIds,omitempty"`
	PublicNetworkAccess                *PublicNetworkAccess            `json:"publicNetworkAccess,omitempty"`
	RestoreParameters                  *RestoreParameters              `json:"restoreParameters,omitempty"`
	VirtualNetworkRules                *[]VirtualNetworkRule           `json:"virtualNetworkRules,omitempty"`
}

var _ json.Unmarshaler = &DatabaseAccountCreateUpdateProperties{}

func (s *DatabaseAccountCreateUpdateProperties) UnmarshalJSON(bytes []byte) error {
	type alias DatabaseAccountCreateUpdateProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into DatabaseAccountCreateUpdateProperties: %+v", err)
	}

	s.AnalyticalStorageConfiguration = decoded.AnalyticalStorageConfiguration
	s.ApiProperties = decoded.ApiProperties
	s.Capabilities = decoded.Capabilities
	s.Capacity = decoded.Capacity
	s.ConnectorOffer = decoded.ConnectorOffer
	s.ConsistencyPolicy = decoded.ConsistencyPolicy
	s.Cors = decoded.Cors
	s.CreateMode = decoded.CreateMode
	s.DatabaseAccountOfferType = decoded.DatabaseAccountOfferType
	s.DefaultIdentity = decoded.DefaultIdentity
	s.DisableKeyBasedMetadataWriteAccess = decoded.DisableKeyBasedMetadataWriteAccess
	s.DisableLocalAuth = decoded.DisableLocalAuth
	s.EnableAnalyticalStorage = decoded.EnableAnalyticalStorage
	s.EnableAutomaticFailover = decoded.EnableAutomaticFailover
	s.EnableCassandraConnector = decoded.EnableCassandraConnector
	s.EnableFreeTier = decoded.EnableFreeTier
	s.EnableMultipleWriteLocations = decoded.EnableMultipleWriteLocations
	s.EnablePartitionMerge = decoded.EnablePartitionMerge
	s.IPRules = decoded.IPRules
	s.IsVirtualNetworkFilterEnabled = decoded.IsVirtualNetworkFilterEnabled
	s.KeyVaultKeyUri = decoded.KeyVaultKeyUri
	s.KeysMetadata = decoded.KeysMetadata
	s.Locations = decoded.Locations
	s.MinimalTlsVersion = decoded.MinimalTlsVersion
	s.NetworkAclBypass = decoded.NetworkAclBypass
	s.NetworkAclBypassResourceIds = decoded.NetworkAclBypassResourceIds
	s.PublicNetworkAccess = decoded.PublicNetworkAccess
	s.RestoreParameters = decoded.RestoreParameters
	s.VirtualNetworkRules = decoded.VirtualNetworkRules

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DatabaseAccountCreateUpdateProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["backupPolicy"]; ok {
		impl, err := unmarshalBackupPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'BackupPolicy' for 'DatabaseAccountCreateUpdateProperties': %+v", err)
		}
		s.BackupPolicy = impl
	}
	return nil
}
