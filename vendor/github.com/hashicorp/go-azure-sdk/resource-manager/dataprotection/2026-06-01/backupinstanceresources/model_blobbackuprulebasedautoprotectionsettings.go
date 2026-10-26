package backupinstanceresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BlobBackupAutoProtectionSettings = BlobBackupRuleBasedAutoProtectionSettings{}

type BlobBackupRuleBasedAutoProtectionSettings struct {
	Rules *[]BlobBackupAutoProtectionRule `json:"rules,omitempty"`

	// Fields inherited from BlobBackupAutoProtectionSettings

	Enabled    bool   `json:"enabled"`
	ObjectType string `json:"objectType"`
}

func (s BlobBackupRuleBasedAutoProtectionSettings) BlobBackupAutoProtectionSettings() BaseBlobBackupAutoProtectionSettingsImpl {
	return BaseBlobBackupAutoProtectionSettingsImpl{
		Enabled:    s.Enabled,
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = BlobBackupRuleBasedAutoProtectionSettings{}

func (s BlobBackupRuleBasedAutoProtectionSettings) MarshalJSON() ([]byte, error) {
	type wrapper BlobBackupRuleBasedAutoProtectionSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobBackupRuleBasedAutoProtectionSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobBackupRuleBasedAutoProtectionSettings: %+v", err)
	}

	decoded["objectType"] = "BlobBackupRuleBasedAutoProtectionSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobBackupRuleBasedAutoProtectionSettings: %+v", err)
	}

	return encoded, nil
}
