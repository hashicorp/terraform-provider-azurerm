package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBlobDatasetTypeProperties struct {
	Compression           *DatasetCompression  `json:"compression,omitempty"`
	FileName              *string              `json:"fileName,omitempty"`
	FolderPath            *string              `json:"folderPath,omitempty"`
	Format                DatasetStorageFormat `json:"format"`
	ModifiedDatetimeEnd   *string              `json:"modifiedDatetimeEnd,omitempty"`
	ModifiedDatetimeStart *string              `json:"modifiedDatetimeStart,omitempty"`
	TableRootLocation     *string              `json:"tableRootLocation,omitempty"`
}

var _ json.Unmarshaler = &AzureBlobDatasetTypeProperties{}

func (s *AzureBlobDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Compression           *DatasetCompression `json:"compression,omitempty"`
		FileName              *string             `json:"fileName,omitempty"`
		FolderPath            *string             `json:"folderPath,omitempty"`
		ModifiedDatetimeEnd   *string             `json:"modifiedDatetimeEnd,omitempty"`
		ModifiedDatetimeStart *string             `json:"modifiedDatetimeStart,omitempty"`
		TableRootLocation     *string             `json:"tableRootLocation,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Compression = decoded.Compression
	s.FileName = decoded.FileName
	s.FolderPath = decoded.FolderPath
	s.ModifiedDatetimeEnd = decoded.ModifiedDatetimeEnd
	s.ModifiedDatetimeStart = decoded.ModifiedDatetimeStart
	s.TableRootLocation = decoded.TableRootLocation

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureBlobDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["format"]; ok {
		impl, err := UnmarshalDatasetStorageFormatImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Format' for 'AzureBlobDatasetTypeProperties': %+v", err)
		}
		s.Format = impl
	}

	return nil
}
