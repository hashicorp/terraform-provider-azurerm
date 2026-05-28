package accountconnectionresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionPropertiesV2 interface {
	ConnectionPropertiesV2() BaseConnectionPropertiesV2Impl
}

var _ ConnectionPropertiesV2 = BaseConnectionPropertiesV2Impl{}

type BaseConnectionPropertiesV2Impl struct {
	AuthType                    ConnectionAuthType    `json:"authType"`
	Category                    *ConnectionCategory   `json:"category,omitempty"`
	CreatedByWorkspaceArmId     *string               `json:"createdByWorkspaceArmId,omitempty"`
	Error                       *string               `json:"error,omitempty"`
	ExpiryTime                  *string               `json:"expiryTime,omitempty"`
	Group                       *ConnectionGroup      `json:"group,omitempty"`
	IsSharedToAll               *bool                 `json:"isSharedToAll,omitempty"`
	Metadata                    *map[string]string    `json:"metadata,omitempty"`
	PeRequirement               *ManagedPERequirement `json:"peRequirement,omitempty"`
	PeStatus                    *ManagedPEStatus      `json:"peStatus,omitempty"`
	SharedUserList              *[]string             `json:"sharedUserList,omitempty"`
	Target                      *string               `json:"target,omitempty"`
	UseWorkspaceManagedIdentity *bool                 `json:"useWorkspaceManagedIdentity,omitempty"`
}

func (s BaseConnectionPropertiesV2Impl) ConnectionPropertiesV2() BaseConnectionPropertiesV2Impl {
	return s
}

var _ ConnectionPropertiesV2 = RawConnectionPropertiesV2Impl{}

// RawConnectionPropertiesV2Impl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawConnectionPropertiesV2Impl struct {
	connectionPropertiesV2 BaseConnectionPropertiesV2Impl
	Type                   string
	Values                 map[string]interface{}
}

func (s RawConnectionPropertiesV2Impl) ConnectionPropertiesV2() BaseConnectionPropertiesV2Impl {
	return s.connectionPropertiesV2
}

func UnmarshalConnectionPropertiesV2Implementation(input []byte) (ConnectionPropertiesV2, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ConnectionPropertiesV2 into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["authType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AAD") {
		var out AADAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AADAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AccessKey") {
		var out AccessKeyAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AccessKeyAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AccountKey") {
		var out AccountKeyAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AccountKeyAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ApiKey") {
		var out ApiKeyAuthConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ApiKeyAuthConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CustomKeys") {
		var out CustomKeysConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomKeysConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ManagedIdentity") {
		var out ManagedIdentityAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedIdentityAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "None") {
		var out NoneAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NoneAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OAuth2") {
		var out OAuth2AuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OAuth2AuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PAT") {
		var out PATAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PATAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SAS") {
		var out SASAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SASAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServicePrincipal") {
		var out ServicePrincipalAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePrincipalAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UsernamePassword") {
		var out UsernamePasswordAuthTypeConnectionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UsernamePasswordAuthTypeConnectionProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseConnectionPropertiesV2Impl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseConnectionPropertiesV2Impl: %+v", err)
	}

	return RawConnectionPropertiesV2Impl{
		connectionPropertiesV2: parent,
		Type:                   value,
		Values:                 temp,
	}, nil

}
