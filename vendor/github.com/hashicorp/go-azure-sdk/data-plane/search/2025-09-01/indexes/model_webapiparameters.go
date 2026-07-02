package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebApiParameters struct {
	AuthIdentity   SearchIndexerDataIdentity `json:"authIdentity"`
	AuthResourceId *string                   `json:"authResourceId,omitempty"`
	HTTPHeaders    *map[string]string        `json:"httpHeaders,omitempty"`
	HTTPMethod     *string                   `json:"httpMethod,omitempty"`
	Timeout        *string                   `json:"timeout,omitempty"`
	Uri            *string                   `json:"uri,omitempty"`
}

var _ json.Unmarshaler = &WebApiParameters{}

func (s *WebApiParameters) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthResourceId *string            `json:"authResourceId,omitempty"`
		HTTPHeaders    *map[string]string `json:"httpHeaders,omitempty"`
		HTTPMethod     *string            `json:"httpMethod,omitempty"`
		Timeout        *string            `json:"timeout,omitempty"`
		Uri            *string            `json:"uri,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthResourceId = decoded.AuthResourceId
	s.HTTPHeaders = decoded.HTTPHeaders
	s.HTTPMethod = decoded.HTTPMethod
	s.Timeout = decoded.Timeout
	s.Uri = decoded.Uri

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling WebApiParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authIdentity"]; ok {
		impl, err := UnmarshalSearchIndexerDataIdentityImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AuthIdentity' for 'WebApiParameters': %+v", err)
		}
		s.AuthIdentity = impl
	}

	return nil
}
