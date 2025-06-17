package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPDatasetTypeProperties struct {
	AdditionalHeaders *interface{}         `json:"additionalHeaders,omitempty"`
	Compression       *DatasetCompression  `json:"compression,omitempty"`
	Format            DatasetStorageFormat `json:"format"`
	RelativeURL       *interface{}         `json:"relativeUrl,omitempty"`
	RequestBody       *interface{}         `json:"requestBody,omitempty"`
	RequestMethod     *interface{}         `json:"requestMethod,omitempty"`
}

var _ json.Unmarshaler = &HTTPDatasetTypeProperties{}

func (s *HTTPDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AdditionalHeaders *interface{}        `json:"additionalHeaders,omitempty"`
		Compression       *DatasetCompression `json:"compression,omitempty"`
		RelativeURL       *interface{}        `json:"relativeUrl,omitempty"`
		RequestBody       *interface{}        `json:"requestBody,omitempty"`
		RequestMethod     *interface{}        `json:"requestMethod,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AdditionalHeaders = decoded.AdditionalHeaders
	s.Compression = decoded.Compression
	s.RelativeURL = decoded.RelativeURL
	s.RequestBody = decoded.RequestBody
	s.RequestMethod = decoded.RequestMethod

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling HTTPDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["format"]; ok {
		impl, err := UnmarshalDatasetStorageFormatImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Format' for 'HTTPDatasetTypeProperties': %+v", err)
		}
		s.Format = impl
	}

	return nil
}
