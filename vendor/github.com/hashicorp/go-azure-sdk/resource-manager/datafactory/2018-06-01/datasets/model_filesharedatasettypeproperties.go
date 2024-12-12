package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileShareDatasetTypeProperties struct {
	Compression           *DatasetCompression  `json:"compression,omitempty"`
	FileFilter            *string              `json:"fileFilter,omitempty"`
	FileName              *string              `json:"fileName,omitempty"`
	FolderPath            *string              `json:"folderPath,omitempty"`
	Format                DatasetStorageFormat `json:"format"`
	ModifiedDatetimeEnd   *string              `json:"modifiedDatetimeEnd,omitempty"`
	ModifiedDatetimeStart *string              `json:"modifiedDatetimeStart,omitempty"`
}

var _ json.Unmarshaler = &FileShareDatasetTypeProperties{}

func (s *FileShareDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Compression           *DatasetCompression `json:"compression,omitempty"`
		FileFilter            *string             `json:"fileFilter,omitempty"`
		FileName              *string             `json:"fileName,omitempty"`
		FolderPath            *string             `json:"folderPath,omitempty"`
		ModifiedDatetimeEnd   *string             `json:"modifiedDatetimeEnd,omitempty"`
		ModifiedDatetimeStart *string             `json:"modifiedDatetimeStart,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Compression = decoded.Compression
	s.FileFilter = decoded.FileFilter
	s.FileName = decoded.FileName
	s.FolderPath = decoded.FolderPath
	s.ModifiedDatetimeEnd = decoded.ModifiedDatetimeEnd
	s.ModifiedDatetimeStart = decoded.ModifiedDatetimeStart

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FileShareDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["format"]; ok {
		impl, err := UnmarshalDatasetStorageFormatImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Format' for 'FileShareDatasetTypeProperties': %+v", err)
		}
		s.Format = impl
	}

	return nil
}
