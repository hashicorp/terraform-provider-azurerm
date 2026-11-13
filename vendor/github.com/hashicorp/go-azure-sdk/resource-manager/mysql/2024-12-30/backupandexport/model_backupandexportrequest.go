package backupandexport

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupAndExportRequest struct {
	BackupSettings BackupSettings     `json:"backupSettings"`
	TargetDetails  BackupStoreDetails `json:"targetDetails"`
}

var _ json.Unmarshaler = &BackupAndExportRequest{}

func (s *BackupAndExportRequest) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		BackupSettings BackupSettings `json:"backupSettings"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.BackupSettings = decoded.BackupSettings

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BackupAndExportRequest into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["targetDetails"]; ok {
		impl, err := UnmarshalBackupStoreDetailsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TargetDetails' for 'BackupAndExportRequest': %+v", err)
		}
		s.TargetDetails = impl
	}

	return nil
}
