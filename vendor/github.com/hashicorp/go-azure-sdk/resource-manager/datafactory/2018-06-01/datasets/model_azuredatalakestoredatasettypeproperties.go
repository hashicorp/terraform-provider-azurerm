package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDataLakeStoreDatasetTypeProperties struct {
	Compression *DatasetCompression  `json:"compression,omitempty"`
	FileName    *string              `json:"fileName,omitempty"`
	FolderPath  *string              `json:"folderPath,omitempty"`
	Format      DatasetStorageFormat `json:"format"`
}

var _ json.Unmarshaler = &AzureDataLakeStoreDatasetTypeProperties{}

func (s *AzureDataLakeStoreDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Compression *DatasetCompression `json:"compression,omitempty"`
		FileName    *string             `json:"fileName,omitempty"`
		FolderPath  *string             `json:"folderPath,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Compression = decoded.Compression
	s.FileName = decoded.FileName
	s.FolderPath = decoded.FolderPath

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureDataLakeStoreDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["format"]; ok {
		impl, err := UnmarshalDatasetStorageFormatImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Format' for 'AzureDataLakeStoreDatasetTypeProperties': %+v", err)
		}
		s.Format = impl
	}

	return nil
}
