package cosmosdb

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseAccountGetProperties struct {
	AnalyticalStorageConfiguration     *AnalyticalStorageConfiguration `json:"analyticalStorageConfiguration,omitempty"`
	ApiProperties                      *ApiProperties                  `json:"apiProperties,omitempty"`
	BackupPolicy                       BackupPolicy                    `json:"backupPolicy"`
	Capabilities                       *[]Capability                   `json:"capabilities,omitempty"`
	Capacity                           *Capacity                       `json:"capacity,omitempty"`
	ConnectorOffer                     *ConnectorOffer                 `json:"connectorOffer,omitempty"`
	ConsistencyPolicy                  *ConsistencyPolicy              `json:"consistencyPolicy,omitempty"`
	Cors                               *[]CorsPolicy                   `json:"cors,omitempty"`
	CreateMode                         *CreateMode                     `json:"createMode,omitempty"`
	DatabaseAccountOfferType           *DatabaseAccountOfferType       `json:"databaseAccountOfferType,omitempty"`
	DefaultIdentity                    *string                         `json:"defaultIdentity,omitempty"`
	DisableKeyBasedMetadataWriteAccess *bool                           `json:"disableKeyBasedMetadataWriteAccess,omitempty"`
	DisableLocalAuth                   *bool                           `json:"disableLocalAuth,omitempty"`
	DocumentEndpoint                   *string                         `json:"documentEndpoint,omitempty"`
	EnableAnalyticalStorage            *bool                           `json:"enableAnalyticalStorage,omitempty"`
	EnableAutomaticFailover            *bool                           `json:"enableAutomaticFailover,omitempty"`
	EnableCassandraConnector           *bool                           `json:"enableCassandraConnector,omitempty"`
	EnableFreeTier                     *bool                           `json:"enableFreeTier,omitempty"`
	EnableMultipleWriteLocations       *bool                           `json:"enableMultipleWriteLocations,omitempty"`
	FailoverPolicies                   *[]FailoverPolicy               `json:"failoverPolicies,omitempty"`
	IPRules                            *[]IPAddressOrRange             `json:"ipRules,omitempty"`
	InstanceId                         *string                         `json:"instanceId,omitempty"`
	IsVirtualNetworkFilterEnabled      *bool                           `json:"isVirtualNetworkFilterEnabled,omitempty"`
	KeyVaultKeyUri                     *string                         `json:"keyVaultKeyUri,omitempty"`
	Locations                          *[]Location                     `json:"locations,omitempty"`
	NetworkAclBypass                   *NetworkAclBypass               `json:"networkAclBypass,omitempty"`
	NetworkAclBypassResourceIds        *[]string                       `json:"networkAclBypassResourceIds,omitempty"`
	PrivateEndpointConnections         *[]PrivateEndpointConnection    `json:"privateEndpointConnections,omitempty"`
	ProvisioningState                  *string                         `json:"provisioningState,omitempty"`
	PublicNetworkAccess                *PublicNetworkAccess            `json:"publicNetworkAccess,omitempty"`
	ReadLocations                      *[]Location                     `json:"readLocations,omitempty"`
	RestoreParameters                  *RestoreParameters              `json:"restoreParameters,omitempty"`
	VirtualNetworkRules                *[]VirtualNetworkRule           `json:"virtualNetworkRules,omitempty"`
	WriteLocations                     *[]Location                     `json:"writeLocations,omitempty"`
}

var _ json.Unmarshaler = &DatabaseAccountGetProperties{}

func (s *DatabaseAccountGetProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AnalyticalStorageConfiguration     *AnalyticalStorageConfiguration `json:"analyticalStorageConfiguration,omitempty"`
		ApiProperties                      *ApiProperties                  `json:"apiProperties,omitempty"`
		Capabilities                       *[]Capability                   `json:"capabilities,omitempty"`
		Capacity                           *Capacity                       `json:"capacity,omitempty"`
		ConnectorOffer                     *ConnectorOffer                 `json:"connectorOffer,omitempty"`
		ConsistencyPolicy                  *ConsistencyPolicy              `json:"consistencyPolicy,omitempty"`
		Cors                               *[]CorsPolicy                   `json:"cors,omitempty"`
		CreateMode                         *CreateMode                     `json:"createMode,omitempty"`
		DatabaseAccountOfferType           *DatabaseAccountOfferType       `json:"databaseAccountOfferType,omitempty"`
		DefaultIdentity                    *string                         `json:"defaultIdentity,omitempty"`
		DisableKeyBasedMetadataWriteAccess *bool                           `json:"disableKeyBasedMetadataWriteAccess,omitempty"`
		DisableLocalAuth                   *bool                           `json:"disableLocalAuth,omitempty"`
		DocumentEndpoint                   *string                         `json:"documentEndpoint,omitempty"`
		EnableAnalyticalStorage            *bool                           `json:"enableAnalyticalStorage,omitempty"`
		EnableAutomaticFailover            *bool                           `json:"enableAutomaticFailover,omitempty"`
		EnableCassandraConnector           *bool                           `json:"enableCassandraConnector,omitempty"`
		EnableFreeTier                     *bool                           `json:"enableFreeTier,omitempty"`
		EnableMultipleWriteLocations       *bool                           `json:"enableMultipleWriteLocations,omitempty"`
		FailoverPolicies                   *[]FailoverPolicy               `json:"failoverPolicies,omitempty"`
		IPRules                            *[]IPAddressOrRange             `json:"ipRules,omitempty"`
		InstanceId                         *string                         `json:"instanceId,omitempty"`
		IsVirtualNetworkFilterEnabled      *bool                           `json:"isVirtualNetworkFilterEnabled,omitempty"`
		KeyVaultKeyUri                     *string                         `json:"keyVaultKeyUri,omitempty"`
		Locations                          *[]Location                     `json:"locations,omitempty"`
		NetworkAclBypass                   *NetworkAclBypass               `json:"networkAclBypass,omitempty"`
		NetworkAclBypassResourceIds        *[]string                       `json:"networkAclBypassResourceIds,omitempty"`
		PrivateEndpointConnections         *[]PrivateEndpointConnection    `json:"privateEndpointConnections,omitempty"`
		ProvisioningState                  *string                         `json:"provisioningState,omitempty"`
		PublicNetworkAccess                *PublicNetworkAccess            `json:"publicNetworkAccess,omitempty"`
		ReadLocations                      *[]Location                     `json:"readLocations,omitempty"`
		RestoreParameters                  *RestoreParameters              `json:"restoreParameters,omitempty"`
		VirtualNetworkRules                *[]VirtualNetworkRule           `json:"virtualNetworkRules,omitempty"`
		WriteLocations                     *[]Location                     `json:"writeLocations,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
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
	s.DocumentEndpoint = decoded.DocumentEndpoint
	s.EnableAnalyticalStorage = decoded.EnableAnalyticalStorage
	s.EnableAutomaticFailover = decoded.EnableAutomaticFailover
	s.EnableCassandraConnector = decoded.EnableCassandraConnector
	s.EnableFreeTier = decoded.EnableFreeTier
	s.EnableMultipleWriteLocations = decoded.EnableMultipleWriteLocations
	s.FailoverPolicies = decoded.FailoverPolicies
	s.IPRules = decoded.IPRules
	s.InstanceId = decoded.InstanceId
	s.IsVirtualNetworkFilterEnabled = decoded.IsVirtualNetworkFilterEnabled
	s.KeyVaultKeyUri = decoded.KeyVaultKeyUri
	s.Locations = decoded.Locations
	s.NetworkAclBypass = decoded.NetworkAclBypass
	s.NetworkAclBypassResourceIds = decoded.NetworkAclBypassResourceIds
	s.PrivateEndpointConnections = decoded.PrivateEndpointConnections
	s.ProvisioningState = decoded.ProvisioningState
	s.PublicNetworkAccess = decoded.PublicNetworkAccess
	s.ReadLocations = decoded.ReadLocations
	s.RestoreParameters = decoded.RestoreParameters
	s.VirtualNetworkRules = decoded.VirtualNetworkRules
	s.WriteLocations = decoded.WriteLocations

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DatabaseAccountGetProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["backupPolicy"]; ok {
		impl, err := UnmarshalBackupPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'BackupPolicy' for 'DatabaseAccountGetProperties': %+v", err)
		}
		s.BackupPolicy = impl
	}

	return nil
}
