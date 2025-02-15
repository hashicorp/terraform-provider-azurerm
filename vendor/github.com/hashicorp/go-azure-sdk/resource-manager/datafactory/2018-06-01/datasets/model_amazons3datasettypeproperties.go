package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmazonS3DatasetTypeProperties struct {
	BucketName            string               `json:"bucketName"`
	Compression           *DatasetCompression  `json:"compression,omitempty"`
	Format                DatasetStorageFormat `json:"format"`
	Key                   *string              `json:"key,omitempty"`
	ModifiedDatetimeEnd   *string              `json:"modifiedDatetimeEnd,omitempty"`
	ModifiedDatetimeStart *string              `json:"modifiedDatetimeStart,omitempty"`
	Prefix                *string              `json:"prefix,omitempty"`
	Version               *string              `json:"version,omitempty"`
}

var _ json.Unmarshaler = &AmazonS3DatasetTypeProperties{}

func (s *AmazonS3DatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		BucketName            string              `json:"bucketName"`
		Compression           *DatasetCompression `json:"compression,omitempty"`
		Key                   *string             `json:"key,omitempty"`
		ModifiedDatetimeEnd   *string             `json:"modifiedDatetimeEnd,omitempty"`
		ModifiedDatetimeStart *string             `json:"modifiedDatetimeStart,omitempty"`
		Prefix                *string             `json:"prefix,omitempty"`
		Version               *string             `json:"version,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.BucketName = decoded.BucketName
	s.Compression = decoded.Compression
	s.Key = decoded.Key
	s.ModifiedDatetimeEnd = decoded.ModifiedDatetimeEnd
	s.ModifiedDatetimeStart = decoded.ModifiedDatetimeStart
	s.Prefix = decoded.Prefix
	s.Version = decoded.Version

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AmazonS3DatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["format"]; ok {
		impl, err := UnmarshalDatasetStorageFormatImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Format' for 'AmazonS3DatasetTypeProperties': %+v", err)
		}
		s.Format = impl
	}

	return nil
}
