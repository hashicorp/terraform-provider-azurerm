package goldengateconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionBaseProperties interface {
	ConnectionBaseProperties() BaseConnectionBasePropertiesImpl
}

var _ ConnectionBaseProperties = BaseConnectionBasePropertiesImpl{}

type BaseConnectionBasePropertiesImpl struct {
	CompartmentId     *string                         `json:"compartmentId,omitempty"`
	ConnectionType    ConnectionType                  `json:"connectionType"`
	DisplayName       string                          `json:"displayName"`
	DoesUseSecretIds  *bool                           `json:"doesUseSecretIds,omitempty"`
	KeyId             *string                         `json:"keyId,omitempty"`
	LifecycleDetails  *string                         `json:"lifecycleDetails,omitempty"`
	LifecycleState    *ConnectionLifecycleState       `json:"lifecycleState,omitempty"`
	NetworkAnchorId   string                          `json:"networkAnchorId"`
	Ocid              *string                         `json:"ocid,omitempty"`
	ProvisioningState *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
	ResourceAnchorId  string                          `json:"resourceAnchorId"`
	RoutingMethod     *RoutingMethod                  `json:"routingMethod,omitempty"`
	TimeCreated       *string                         `json:"timeCreated,omitempty"`
	TimeUpdated       *string                         `json:"timeUpdated,omitempty"`
	VaultId           *string                         `json:"vaultId,omitempty"`
}

func (s BaseConnectionBasePropertiesImpl) ConnectionBaseProperties() BaseConnectionBasePropertiesImpl {
	return s
}

var _ ConnectionBaseProperties = RawConnectionBasePropertiesImpl{}

// RawConnectionBasePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawConnectionBasePropertiesImpl struct {
	connectionBaseProperties BaseConnectionBasePropertiesImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawConnectionBasePropertiesImpl) ConnectionBaseProperties() BaseConnectionBasePropertiesImpl {
	return s.connectionBaseProperties
}

func (s RawConnectionBasePropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalConnectionBasePropertiesImplementation(input []byte) (ConnectionBaseProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectionBaseProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["connectionType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "KAFKA") {
		var out KafkaConnectionDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KafkaConnectionDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MICROSOFT_FABRIC") {
		var out MicrosoftFabricConnectionDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftFabricConnectionDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ORACLE") {
		var out OracleConnectionDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleConnectionDetails: %+v", err)
		}
		return out, nil
	}

	var parent BaseConnectionBasePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseConnectionBasePropertiesImpl: %+v", err)
	}

	return RawConnectionBasePropertiesImpl{
		connectionBaseProperties: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
