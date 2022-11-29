package servers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ServerPropertiesForCreate = ServerPropertiesForDefaultCreate{}

type ServerPropertiesForDefaultCreate struct {
	AdministratorLogin         string `json:"administratorLogin"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`

	// Fields inherited from ServerPropertiesForCreate
	MinimalTlsVersion   *MinimalTlsVersionEnum   `json:"minimalTlsVersion,omitempty"`
	PublicNetworkAccess *PublicNetworkAccessEnum `json:"publicNetworkAccess,omitempty"`
	SslEnforcement      *SslEnforcementEnum      `json:"sslEnforcement,omitempty"`
	StorageProfile      *StorageProfile          `json:"storageProfile,omitempty"`
	Version             *ServerVersion           `json:"version,omitempty"`
}

var _ json.Marshaler = ServerPropertiesForDefaultCreate{}

func (s ServerPropertiesForDefaultCreate) MarshalJSON() ([]byte, error) {
	type wrapper ServerPropertiesForDefaultCreate
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServerPropertiesForDefaultCreate: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServerPropertiesForDefaultCreate: %+v", err)
	}
	decoded["createMode"] = "Default"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServerPropertiesForDefaultCreate: %+v", err)
	}

	return encoded, nil
}
