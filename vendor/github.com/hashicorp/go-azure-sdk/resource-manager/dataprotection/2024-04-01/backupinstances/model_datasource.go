package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Datasource struct {
	DatasourceType     *string                `json:"datasourceType,omitempty"`
	ObjectType         *string                `json:"objectType,omitempty"`
	ResourceID         string                 `json:"resourceID"`
	ResourceLocation   *string                `json:"resourceLocation,omitempty"`
	ResourceName       *string                `json:"resourceName,omitempty"`
	ResourceProperties BaseResourceProperties `json:"resourceProperties"`
	ResourceType       *string                `json:"resourceType,omitempty"`
	ResourceUri        *string                `json:"resourceUri,omitempty"`
}

var _ json.Unmarshaler = &Datasource{}

func (s *Datasource) UnmarshalJSON(bytes []byte) error {
	type alias Datasource
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into Datasource: %+v", err)
	}

	s.DatasourceType = decoded.DatasourceType
	s.ObjectType = decoded.ObjectType
	s.ResourceID = decoded.ResourceID
	s.ResourceLocation = decoded.ResourceLocation
	s.ResourceName = decoded.ResourceName
	s.ResourceType = decoded.ResourceType
	s.ResourceUri = decoded.ResourceUri

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Datasource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["resourceProperties"]; ok {
		impl, err := unmarshalBaseResourcePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ResourceProperties' for 'Datasource': %+v", err)
		}
		s.ResourceProperties = impl
	}
	return nil
}
