package settings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Setting = DataExportSettings{}

type DataExportSettings struct {
	Properties *DataExportSettingProperties `json:"properties,omitempty"`

	// Fields inherited from Setting

	Id   *string     `json:"id,omitempty"`
	Kind SettingKind `json:"kind"`
	Name *string     `json:"name,omitempty"`
	Type *string     `json:"type,omitempty"`
}

func (s DataExportSettings) Setting() BaseSettingImpl {
	return BaseSettingImpl{
		Id:   s.Id,
		Kind: s.Kind,
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = DataExportSettings{}

func (s DataExportSettings) MarshalJSON() ([]byte, error) {
	type wrapper DataExportSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DataExportSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DataExportSettings: %+v", err)
	}

	decoded["kind"] = "DataExportSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DataExportSettings: %+v", err)
	}

	return encoded, nil
}
