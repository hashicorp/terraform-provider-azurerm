package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConfiguration struct {
	TransportFileShareConfiguration FileShareConfiguration `json:"transportFileShareConfiguration"`
}

var _ json.Unmarshaler = &StorageConfiguration{}

func (s *StorageConfiguration) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling StorageConfiguration into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["transportFileShareConfiguration"]; ok {
		impl, err := UnmarshalFileShareConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TransportFileShareConfiguration' for 'StorageConfiguration': %+v", err)
		}
		s.TransportFileShareConfiguration = impl
	}

	return nil
}
