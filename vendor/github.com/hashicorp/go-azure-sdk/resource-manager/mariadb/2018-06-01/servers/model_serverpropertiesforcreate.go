package servers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPropertiesForCreate interface {
	ServerPropertiesForCreate() BaseServerPropertiesForCreateImpl
}

var _ ServerPropertiesForCreate = BaseServerPropertiesForCreateImpl{}

type BaseServerPropertiesForCreateImpl struct {
	CreateMode          CreateMode               `json:"createMode"`
	MinimalTlsVersion   *MinimalTlsVersionEnum   `json:"minimalTlsVersion,omitempty"`
	PublicNetworkAccess *PublicNetworkAccessEnum `json:"publicNetworkAccess,omitempty"`
	SslEnforcement      *SslEnforcementEnum      `json:"sslEnforcement,omitempty"`
	StorageProfile      *StorageProfile          `json:"storageProfile,omitempty"`
	Version             *ServerVersion           `json:"version,omitempty"`
}

func (s BaseServerPropertiesForCreateImpl) ServerPropertiesForCreate() BaseServerPropertiesForCreateImpl {
	return s
}

var _ ServerPropertiesForCreate = RawServerPropertiesForCreateImpl{}

// RawServerPropertiesForCreateImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawServerPropertiesForCreateImpl struct {
	serverPropertiesForCreate BaseServerPropertiesForCreateImpl
	Type                      string
	Values                    map[string]interface{}
}

func (s RawServerPropertiesForCreateImpl) ServerPropertiesForCreate() BaseServerPropertiesForCreateImpl {
	return s.serverPropertiesForCreate
}

func UnmarshalServerPropertiesForCreateImplementation(input []byte) (ServerPropertiesForCreate, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ServerPropertiesForCreate into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["createMode"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Default") {
		var out ServerPropertiesForDefaultCreate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPropertiesForDefaultCreate: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GeoRestore") {
		var out ServerPropertiesForGeoRestore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPropertiesForGeoRestore: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Replica") {
		var out ServerPropertiesForReplica
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPropertiesForReplica: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PointInTimeRestore") {
		var out ServerPropertiesForRestore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPropertiesForRestore: %+v", err)
		}
		return out, nil
	}

	var parent BaseServerPropertiesForCreateImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseServerPropertiesForCreateImpl: %+v", err)
	}

	return RawServerPropertiesForCreateImpl{
		serverPropertiesForCreate: parent,
		Type:                      value,
		Values:                    temp,
	}, nil

}
