package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StoreReadSettings interface {
	StoreReadSettings() BaseStoreReadSettingsImpl
}

var _ StoreReadSettings = BaseStoreReadSettingsImpl{}

type BaseStoreReadSettingsImpl struct {
	DisableMetricsCollection *bool  `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64 `json:"maxConcurrentConnections,omitempty"`
	Type                     string `json:"type"`
}

func (s BaseStoreReadSettingsImpl) StoreReadSettings() BaseStoreReadSettingsImpl {
	return s
}

var _ StoreReadSettings = RawStoreReadSettingsImpl{}

// RawStoreReadSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStoreReadSettingsImpl struct {
	storeReadSettings BaseStoreReadSettingsImpl
	Type              string
	Values            map[string]interface{}
}

func (s RawStoreReadSettingsImpl) StoreReadSettings() BaseStoreReadSettingsImpl {
	return s.storeReadSettings
}

func UnmarshalStoreReadSettingsImplementation(input []byte) (StoreReadSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StoreReadSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AmazonS3CompatibleReadSettings") {
		var out AmazonS3CompatibleReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonS3CompatibleReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonS3ReadSettings") {
		var out AmazonS3ReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonS3ReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobFSReadSettings") {
		var out AzureBlobFSReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobFSReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobStorageReadSettings") {
		var out AzureBlobStorageReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobStorageReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeStoreReadSettings") {
		var out AzureDataLakeStoreReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeStoreReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureFileStorageReadSettings") {
		var out AzureFileStorageReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileStorageReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileServerReadSettings") {
		var out FileServerReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileServerReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FtpReadSettings") {
		var out FtpReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FtpReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleCloudStorageReadSettings") {
		var out GoogleCloudStorageReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleCloudStorageReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HttpReadSettings") {
		var out HTTPReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HTTPReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HdfsReadSettings") {
		var out HdfsReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HdfsReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LakeHouseReadSettings") {
		var out LakeHouseReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LakeHouseReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleCloudStorageReadSettings") {
		var out OracleCloudStorageReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleCloudStorageReadSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SftpReadSettings") {
		var out SftpReadSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SftpReadSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseStoreReadSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStoreReadSettingsImpl: %+v", err)
	}

	return RawStoreReadSettingsImpl{
		storeReadSettings: parent,
		Type:              value,
		Values:            temp,
	}, nil

}
