package datasets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatasetLocation interface {
	DatasetLocation() BaseDatasetLocationImpl
}

var _ DatasetLocation = BaseDatasetLocationImpl{}

type BaseDatasetLocationImpl struct {
	FileName   *string `json:"fileName,omitempty"`
	FolderPath *string `json:"folderPath,omitempty"`
	Type       string  `json:"type"`
}

func (s BaseDatasetLocationImpl) DatasetLocation() BaseDatasetLocationImpl {
	return s
}

var _ DatasetLocation = RawDatasetLocationImpl{}

// RawDatasetLocationImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDatasetLocationImpl struct {
	datasetLocation BaseDatasetLocationImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawDatasetLocationImpl) DatasetLocation() BaseDatasetLocationImpl {
	return s.datasetLocation
}

func UnmarshalDatasetLocationImplementation(input []byte) (DatasetLocation, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DatasetLocation into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AmazonS3CompatibleLocation") {
		var out AmazonS3CompatibleLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonS3CompatibleLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonS3Location") {
		var out AmazonS3Location
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonS3Location: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobFSLocation") {
		var out AzureBlobFSLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobFSLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobStorageLocation") {
		var out AzureBlobStorageLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobStorageLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeStoreLocation") {
		var out AzureDataLakeStoreLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeStoreLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureFileStorageLocation") {
		var out AzureFileStorageLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileStorageLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileServerLocation") {
		var out FileServerLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileServerLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FtpServerLocation") {
		var out FtpServerLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FtpServerLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleCloudStorageLocation") {
		var out GoogleCloudStorageLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleCloudStorageLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HttpServerLocation") {
		var out HTTPServerLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HTTPServerLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HdfsLocation") {
		var out HdfsLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HdfsLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LakeHouseLocation") {
		var out LakeHouseLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LakeHouseLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleCloudStorageLocation") {
		var out OracleCloudStorageLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleCloudStorageLocation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SftpLocation") {
		var out SftpLocation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SftpLocation: %+v", err)
		}
		return out, nil
	}

	var parent BaseDatasetLocationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDatasetLocationImpl: %+v", err)
	}

	return RawDatasetLocationImpl{
		datasetLocation: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
