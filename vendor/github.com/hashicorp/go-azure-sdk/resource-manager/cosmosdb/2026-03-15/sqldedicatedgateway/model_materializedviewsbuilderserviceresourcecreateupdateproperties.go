package sqldedicatedgateway

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ServiceResourceCreateUpdateProperties = MaterializedViewsBuilderServiceResourceCreateUpdateProperties{}

type MaterializedViewsBuilderServiceResourceCreateUpdateProperties struct {

	// Fields inherited from ServiceResourceCreateUpdateProperties

	InstanceCount *int64       `json:"instanceCount,omitempty"`
	InstanceSize  *ServiceSize `json:"instanceSize,omitempty"`
	ServiceType   ServiceType  `json:"serviceType"`
}

func (s MaterializedViewsBuilderServiceResourceCreateUpdateProperties) ServiceResourceCreateUpdateProperties() BaseServiceResourceCreateUpdatePropertiesImpl {
	return BaseServiceResourceCreateUpdatePropertiesImpl{
		InstanceCount: s.InstanceCount,
		InstanceSize:  s.InstanceSize,
		ServiceType:   s.ServiceType,
	}
}

var _ json.Marshaler = MaterializedViewsBuilderServiceResourceCreateUpdateProperties{}

func (s MaterializedViewsBuilderServiceResourceCreateUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper MaterializedViewsBuilderServiceResourceCreateUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MaterializedViewsBuilderServiceResourceCreateUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MaterializedViewsBuilderServiceResourceCreateUpdateProperties: %+v", err)
	}

	decoded["serviceType"] = "MaterializedViewsBuilder"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MaterializedViewsBuilderServiceResourceCreateUpdateProperties: %+v", err)
	}

	return encoded, nil
}
