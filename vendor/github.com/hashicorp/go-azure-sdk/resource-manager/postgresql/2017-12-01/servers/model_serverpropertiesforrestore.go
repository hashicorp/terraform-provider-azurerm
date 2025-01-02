package servers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ServerPropertiesForCreate = ServerPropertiesForRestore{}

type ServerPropertiesForRestore struct {
	RestorePointInTime string `json:"restorePointInTime"`
	SourceServerId     string `json:"sourceServerId"`

	// Fields inherited from ServerPropertiesForCreate

	CreateMode               CreateMode                `json:"createMode"`
	InfrastructureEncryption *InfrastructureEncryption `json:"infrastructureEncryption,omitempty"`
	MinimalTlsVersion        *MinimalTlsVersionEnum    `json:"minimalTlsVersion,omitempty"`
	PublicNetworkAccess      *PublicNetworkAccessEnum  `json:"publicNetworkAccess,omitempty"`
	SslEnforcement           *SslEnforcementEnum       `json:"sslEnforcement,omitempty"`
	StorageProfile           *StorageProfile           `json:"storageProfile,omitempty"`
	Version                  *ServerVersion            `json:"version,omitempty"`
}

func (s ServerPropertiesForRestore) ServerPropertiesForCreate() BaseServerPropertiesForCreateImpl {
	return BaseServerPropertiesForCreateImpl{
		CreateMode:               s.CreateMode,
		InfrastructureEncryption: s.InfrastructureEncryption,
		MinimalTlsVersion:        s.MinimalTlsVersion,
		PublicNetworkAccess:      s.PublicNetworkAccess,
		SslEnforcement:           s.SslEnforcement,
		StorageProfile:           s.StorageProfile,
		Version:                  s.Version,
	}
}

var _ json.Marshaler = ServerPropertiesForRestore{}

func (s ServerPropertiesForRestore) MarshalJSON() ([]byte, error) {
	type wrapper ServerPropertiesForRestore
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServerPropertiesForRestore: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServerPropertiesForRestore: %+v", err)
	}

	decoded["createMode"] = "PointInTimeRestore"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServerPropertiesForRestore: %+v", err)
	}

	return encoded, nil
}
