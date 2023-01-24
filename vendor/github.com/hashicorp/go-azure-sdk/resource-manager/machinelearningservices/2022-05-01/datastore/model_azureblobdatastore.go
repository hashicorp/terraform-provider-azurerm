package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Datastore = AzureBlobDatastore{}

type AzureBlobDatastore struct {
	AccountName                   *string                        `json:"accountName,omitempty"`
	ContainerName                 *string                        `json:"containerName,omitempty"`
	Endpoint                      *string                        `json:"endpoint,omitempty"`
	Protocol                      *string                        `json:"protocol,omitempty"`
	ServiceDataAccessAuthIdentity *ServiceDataAccessAuthIdentity `json:"serviceDataAccessAuthIdentity,omitempty"`

	// Fields inherited from Datastore
	Credentials DatastoreCredentials `json:"credentials"`
	Description *string              `json:"description,omitempty"`
	IsDefault   *bool                `json:"isDefault,omitempty"`
	Properties  *map[string]string   `json:"properties,omitempty"`
	Tags        *map[string]string   `json:"tags,omitempty"`
}

var _ json.Marshaler = AzureBlobDatastore{}

func (s AzureBlobDatastore) MarshalJSON() ([]byte, error) {
	type wrapper AzureBlobDatastore
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureBlobDatastore: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBlobDatastore: %+v", err)
	}
	decoded["datastoreType"] = "AzureBlob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureBlobDatastore: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureBlobDatastore{}

func (s *AzureBlobDatastore) UnmarshalJSON(bytes []byte) error {
	type alias AzureBlobDatastore
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into AzureBlobDatastore: %+v", err)
	}

	s.AccountName = decoded.AccountName
	s.ContainerName = decoded.ContainerName
	s.Description = decoded.Description
	s.Endpoint = decoded.Endpoint
	s.IsDefault = decoded.IsDefault
	s.Properties = decoded.Properties
	s.Protocol = decoded.Protocol
	s.ServiceDataAccessAuthIdentity = decoded.ServiceDataAccessAuthIdentity
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureBlobDatastore into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["credentials"]; ok {
		impl, err := unmarshalDatastoreCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credentials' for 'AzureBlobDatastore': %+v", err)
		}
		s.Credentials = impl
	}
	return nil
}
