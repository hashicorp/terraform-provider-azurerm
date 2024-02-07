package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateRestoreRequestObject struct {
	RestoreRequestObject AzureBackupRestoreRequest `json:"restoreRequestObject"`
}

var _ json.Unmarshaler = &ValidateRestoreRequestObject{}

func (s *ValidateRestoreRequestObject) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ValidateRestoreRequestObject into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["restoreRequestObject"]; ok {
		impl, err := unmarshalAzureBackupRestoreRequestImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RestoreRequestObject' for 'ValidateRestoreRequestObject': %+v", err)
		}
		s.RestoreRequestObject = impl
	}
	return nil
}
