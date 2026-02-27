package sqldedicatedgateway

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ServiceResourceCreateUpdateProperties = DataTransferServiceResourceCreateUpdateProperties{}

type DataTransferServiceResourceCreateUpdateProperties struct {

	// Fields inherited from ServiceResourceCreateUpdateProperties

	InstanceCount *int64       `json:"instanceCount,omitempty"`
	InstanceSize  *ServiceSize `json:"instanceSize,omitempty"`
	ServiceType   ServiceType  `json:"serviceType"`
}

func (s DataTransferServiceResourceCreateUpdateProperties) ServiceResourceCreateUpdateProperties() BaseServiceResourceCreateUpdatePropertiesImpl {
	return BaseServiceResourceCreateUpdatePropertiesImpl{
		InstanceCount: s.InstanceCount,
		InstanceSize:  s.InstanceSize,
		ServiceType:   s.ServiceType,
	}
}

var _ json.Marshaler = DataTransferServiceResourceCreateUpdateProperties{}

func (s DataTransferServiceResourceCreateUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper DataTransferServiceResourceCreateUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DataTransferServiceResourceCreateUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DataTransferServiceResourceCreateUpdateProperties: %+v", err)
	}

	decoded["serviceType"] = "DataTransfer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DataTransferServiceResourceCreateUpdateProperties: %+v", err)
	}

	return encoded, nil
}
