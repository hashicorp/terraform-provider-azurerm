package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectorResourceProperties struct {
	ConnectorBasicInfo       *ConnectorInfoBase           `json:"connectorBasicInfo,omitempty"`
	ConnectorServiceTypeInfo ConnectorServiceTypeInfoBase `json:"connectorServiceTypeInfo"`
	PartnerConnectorInfo     PartnerInfoBase              `json:"partnerConnectorInfo"`
}

var _ json.Unmarshaler = &ConnectorResourceProperties{}

func (s *ConnectorResourceProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ConnectorBasicInfo *ConnectorInfoBase `json:"connectorBasicInfo,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ConnectorBasicInfo = decoded.ConnectorBasicInfo

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ConnectorResourceProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["connectorServiceTypeInfo"]; ok {
		impl, err := UnmarshalConnectorServiceTypeInfoBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ConnectorServiceTypeInfo' for 'ConnectorResourceProperties': %+v", err)
		}
		s.ConnectorServiceTypeInfo = impl
	}

	if v, ok := temp["partnerConnectorInfo"]; ok {
		impl, err := UnmarshalPartnerInfoBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'PartnerConnectorInfo' for 'ConnectorResourceProperties': %+v", err)
		}
		s.PartnerConnectorInfo = impl
	}

	return nil
}
