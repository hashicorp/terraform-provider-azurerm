package datasources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerDataSource struct {
	Container                   SearchIndexerDataContainer   `json:"container"`
	Credentials                 DataSourceCredentials        `json:"credentials"`
	DataChangeDetectionPolicy   DataChangeDetectionPolicy    `json:"dataChangeDetectionPolicy"`
	DataDeletionDetectionPolicy DataDeletionDetectionPolicy  `json:"dataDeletionDetectionPolicy"`
	Description                 *string                      `json:"description,omitempty"`
	EncryptionKey               *SearchResourceEncryptionKey `json:"encryptionKey,omitempty"`
	Name                        string                       `json:"name"`
	OdataEtag                   *string                      `json:"@odata.etag,omitempty"`
	Type                        SearchIndexerDataSourceType  `json:"type"`
}

var _ json.Unmarshaler = &SearchIndexerDataSource{}

func (s *SearchIndexerDataSource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Container     SearchIndexerDataContainer   `json:"container"`
		Credentials   DataSourceCredentials        `json:"credentials"`
		Description   *string                      `json:"description,omitempty"`
		EncryptionKey *SearchResourceEncryptionKey `json:"encryptionKey,omitempty"`
		Name          string                       `json:"name"`
		OdataEtag     *string                      `json:"@odata.etag,omitempty"`
		Type          SearchIndexerDataSourceType  `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Container = decoded.Container
	s.Credentials = decoded.Credentials
	s.Description = decoded.Description
	s.EncryptionKey = decoded.EncryptionKey
	s.Name = decoded.Name
	s.OdataEtag = decoded.OdataEtag
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SearchIndexerDataSource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["dataChangeDetectionPolicy"]; ok {
		impl, err := UnmarshalDataChangeDetectionPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'DataChangeDetectionPolicy' for 'SearchIndexerDataSource': %+v", err)
		}
		s.DataChangeDetectionPolicy = impl
	}

	if v, ok := temp["dataDeletionDetectionPolicy"]; ok {
		impl, err := UnmarshalDataDeletionDetectionPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'DataDeletionDetectionPolicy' for 'SearchIndexerDataSource': %+v", err)
		}
		s.DataDeletionDetectionPolicy = impl
	}

	return nil
}
