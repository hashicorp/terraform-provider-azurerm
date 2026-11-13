package backupinstanceresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobBackupAutoProtectionSettings interface {
	BlobBackupAutoProtectionSettings() BaseBlobBackupAutoProtectionSettingsImpl
}

var _ BlobBackupAutoProtectionSettings = BaseBlobBackupAutoProtectionSettingsImpl{}

type BaseBlobBackupAutoProtectionSettingsImpl struct {
	Enabled    bool   `json:"enabled"`
	ObjectType string `json:"objectType"`
}

func (s BaseBlobBackupAutoProtectionSettingsImpl) BlobBackupAutoProtectionSettings() BaseBlobBackupAutoProtectionSettingsImpl {
	return s
}

var _ BlobBackupAutoProtectionSettings = RawBlobBackupAutoProtectionSettingsImpl{}

// RawBlobBackupAutoProtectionSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawBlobBackupAutoProtectionSettingsImpl struct {
	blobBackupAutoProtectionSettings BaseBlobBackupAutoProtectionSettingsImpl
	Type                             string
	Values                           map[string]interface{}
}

func (s RawBlobBackupAutoProtectionSettingsImpl) BlobBackupAutoProtectionSettings() BaseBlobBackupAutoProtectionSettingsImpl {
	return s.blobBackupAutoProtectionSettings
}

func (s RawBlobBackupAutoProtectionSettingsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalBlobBackupAutoProtectionSettingsImplementation(input []byte) (BlobBackupAutoProtectionSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobBackupAutoProtectionSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BlobBackupRuleBasedAutoProtectionSettings") {
		var out BlobBackupRuleBasedAutoProtectionSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobBackupRuleBasedAutoProtectionSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseBlobBackupAutoProtectionSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBlobBackupAutoProtectionSettingsImpl: %+v", err)
	}

	return RawBlobBackupAutoProtectionSettingsImpl{
		blobBackupAutoProtectionSettings: parent,
		Type:                             value,
		Values:                           temp,
	}, nil

}
