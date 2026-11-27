package storagecontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageContainerProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	StorageStore      StorageStore       `json:"storageStore"`
}

var _ json.Unmarshaler = &StorageContainerProperties{}

func (s *StorageContainerProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ProvisioningState = decoded.ProvisioningState

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling StorageContainerProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["storageStore"]; ok {
		impl, err := UnmarshalStorageStoreImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'StorageStore' for 'StorageContainerProperties': %+v", err)
		}
		s.StorageStore = impl
	}

	return nil
}
