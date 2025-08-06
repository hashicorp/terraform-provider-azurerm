package linkedservices

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebLinkedServiceTypeProperties interface {
	WebLinkedServiceTypeProperties() BaseWebLinkedServiceTypePropertiesImpl
}

var _ WebLinkedServiceTypeProperties = BaseWebLinkedServiceTypePropertiesImpl{}

type BaseWebLinkedServiceTypePropertiesImpl struct {
	AuthenticationType WebAuthenticationType `json:"authenticationType"`
	Url                interface{}           `json:"url"`
}

func (s BaseWebLinkedServiceTypePropertiesImpl) WebLinkedServiceTypeProperties() BaseWebLinkedServiceTypePropertiesImpl {
	return s
}

var _ WebLinkedServiceTypeProperties = RawWebLinkedServiceTypePropertiesImpl{}

// RawWebLinkedServiceTypePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawWebLinkedServiceTypePropertiesImpl struct {
	webLinkedServiceTypeProperties BaseWebLinkedServiceTypePropertiesImpl
	Type                           string
	Values                         map[string]interface{}
}

func (s RawWebLinkedServiceTypePropertiesImpl) WebLinkedServiceTypeProperties() BaseWebLinkedServiceTypePropertiesImpl {
	return s.webLinkedServiceTypeProperties
}

func UnmarshalWebLinkedServiceTypePropertiesImplementation(input []byte) (WebLinkedServiceTypeProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling WebLinkedServiceTypeProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["authenticationType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Anonymous") {
		var out WebAnonymousAuthentication
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebAnonymousAuthentication: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Basic") {
		var out WebBasicAuthentication
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebBasicAuthentication: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ClientCertificate") {
		var out WebClientCertificateAuthentication
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebClientCertificateAuthentication: %+v", err)
		}
		return out, nil
	}

	var parent BaseWebLinkedServiceTypePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseWebLinkedServiceTypePropertiesImpl: %+v", err)
	}

	return RawWebLinkedServiceTypePropertiesImpl{
		webLinkedServiceTypeProperties: parent,
		Type:                           value,
		Values:                         temp,
	}, nil

}
