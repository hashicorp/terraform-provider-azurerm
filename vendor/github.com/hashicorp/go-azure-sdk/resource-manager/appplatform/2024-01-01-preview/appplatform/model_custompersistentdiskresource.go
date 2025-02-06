package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomPersistentDiskResource struct {
	CustomPersistentDiskProperties CustomPersistentDiskProperties `json:"customPersistentDiskProperties"`
	StorageId                      string                         `json:"storageId"`
}

var _ json.Unmarshaler = &CustomPersistentDiskResource{}

func (s *CustomPersistentDiskResource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		StorageId string `json:"storageId"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.StorageId = decoded.StorageId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CustomPersistentDiskResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customPersistentDiskProperties"]; ok {
		impl, err := UnmarshalCustomPersistentDiskPropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomPersistentDiskProperties' for 'CustomPersistentDiskResource': %+v", err)
		}
		s.CustomPersistentDiskProperties = impl
	}

	return nil
}
