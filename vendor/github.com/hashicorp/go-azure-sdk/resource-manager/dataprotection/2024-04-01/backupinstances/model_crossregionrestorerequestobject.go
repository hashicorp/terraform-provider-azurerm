package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CrossRegionRestoreRequestObject struct {
	CrossRegionRestoreDetails CrossRegionRestoreDetails `json:"crossRegionRestoreDetails"`
	RestoreRequestObject      AzureBackupRestoreRequest `json:"restoreRequestObject"`
}

var _ json.Unmarshaler = &CrossRegionRestoreRequestObject{}

func (s *CrossRegionRestoreRequestObject) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CrossRegionRestoreDetails CrossRegionRestoreDetails `json:"crossRegionRestoreDetails"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CrossRegionRestoreDetails = decoded.CrossRegionRestoreDetails

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CrossRegionRestoreRequestObject into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["restoreRequestObject"]; ok {
		impl, err := UnmarshalAzureBackupRestoreRequestImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RestoreRequestObject' for 'CrossRegionRestoreRequestObject': %+v", err)
		}
		s.RestoreRequestObject = impl
	}

	return nil
}
