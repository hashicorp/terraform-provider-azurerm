package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BaseResourceProperties = DefaultResourceProperties{}

type DefaultResourceProperties struct {

	// Fields inherited from BaseResourceProperties

	ObjectType ResourcePropertiesObjectType `json:"objectType"`
}

func (s DefaultResourceProperties) BaseResourceProperties() BaseBaseResourcePropertiesImpl {
	return BaseBaseResourcePropertiesImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = DefaultResourceProperties{}

func (s DefaultResourceProperties) MarshalJSON() ([]byte, error) {
	type wrapper DefaultResourceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DefaultResourceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DefaultResourceProperties: %+v", err)
	}

	decoded["objectType"] = "DefaultResourceProperties"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DefaultResourceProperties: %+v", err)
	}

	return encoded, nil
}
