package sqldedicatedgateway

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ServiceResourceCreateUpdateProperties = SqlDedicatedGatewayServiceResourceCreateUpdateProperties{}

type SqlDedicatedGatewayServiceResourceCreateUpdateProperties struct {
	DedicatedGatewayType *DedicatedGatewayType `json:"dedicatedGatewayType,omitempty"`

	// Fields inherited from ServiceResourceCreateUpdateProperties

	InstanceCount *int64       `json:"instanceCount,omitempty"`
	InstanceSize  *ServiceSize `json:"instanceSize,omitempty"`
	ServiceType   ServiceType  `json:"serviceType"`
}

func (s SqlDedicatedGatewayServiceResourceCreateUpdateProperties) ServiceResourceCreateUpdateProperties() BaseServiceResourceCreateUpdatePropertiesImpl {
	return BaseServiceResourceCreateUpdatePropertiesImpl{
		InstanceCount: s.InstanceCount,
		InstanceSize:  s.InstanceSize,
		ServiceType:   s.ServiceType,
	}
}

var _ json.Marshaler = SqlDedicatedGatewayServiceResourceCreateUpdateProperties{}

func (s SqlDedicatedGatewayServiceResourceCreateUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper SqlDedicatedGatewayServiceResourceCreateUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SqlDedicatedGatewayServiceResourceCreateUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SqlDedicatedGatewayServiceResourceCreateUpdateProperties: %+v", err)
	}

	decoded["serviceType"] = "SqlDedicatedGateway"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SqlDedicatedGatewayServiceResourceCreateUpdateProperties: %+v", err)
	}

	return encoded, nil
}
