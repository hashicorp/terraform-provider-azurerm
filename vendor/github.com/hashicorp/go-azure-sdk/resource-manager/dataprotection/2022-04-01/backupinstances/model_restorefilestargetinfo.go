package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RestoreTargetInfoBase = RestoreFilesTargetInfo{}

type RestoreFilesTargetInfo struct {
	TargetDetails TargetDetails `json:"targetDetails"`

	// Fields inherited from RestoreTargetInfoBase
	RecoveryOption  RecoveryOption `json:"recoveryOption"`
	RestoreLocation *string        `json:"restoreLocation,omitempty"`
}

var _ json.Marshaler = RestoreFilesTargetInfo{}

func (s RestoreFilesTargetInfo) MarshalJSON() ([]byte, error) {
	type wrapper RestoreFilesTargetInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RestoreFilesTargetInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RestoreFilesTargetInfo: %+v", err)
	}
	decoded["objectType"] = "RestoreFilesTargetInfo"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RestoreFilesTargetInfo: %+v", err)
	}

	return encoded, nil
}
