package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StoreWriteSettings interface {
	StoreWriteSettings() BaseStoreWriteSettingsImpl
}

var _ StoreWriteSettings = BaseStoreWriteSettingsImpl{}

type BaseStoreWriteSettingsImpl struct {
	CopyBehavior             *string         `json:"copyBehavior,omitempty"`
	DisableMetricsCollection *bool           `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64          `json:"maxConcurrentConnections,omitempty"`
	Metadata                 *[]MetadataItem `json:"metadata,omitempty"`
	Type                     string          `json:"type"`
}

func (s BaseStoreWriteSettingsImpl) StoreWriteSettings() BaseStoreWriteSettingsImpl {
	return s
}

var _ StoreWriteSettings = RawStoreWriteSettingsImpl{}

// RawStoreWriteSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStoreWriteSettingsImpl struct {
	storeWriteSettings BaseStoreWriteSettingsImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawStoreWriteSettingsImpl) StoreWriteSettings() BaseStoreWriteSettingsImpl {
	return s.storeWriteSettings
}

func UnmarshalStoreWriteSettingsImplementation(input []byte) (StoreWriteSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StoreWriteSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureBlobFSWriteSettings") {
		var out AzureBlobFSWriteSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobFSWriteSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobStorageWriteSettings") {
		var out AzureBlobStorageWriteSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobStorageWriteSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeStoreWriteSettings") {
		var out AzureDataLakeStoreWriteSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeStoreWriteSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureFileStorageWriteSettings") {
		var out AzureFileStorageWriteSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileStorageWriteSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileServerWriteSettings") {
		var out FileServerWriteSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileServerWriteSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LakeHouseWriteSettings") {
		var out LakeHouseWriteSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LakeHouseWriteSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SftpWriteSettings") {
		var out SftpWriteSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SftpWriteSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseStoreWriteSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStoreWriteSettingsImpl: %+v", err)
	}

	return RawStoreWriteSettingsImpl{
		storeWriteSettings: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
