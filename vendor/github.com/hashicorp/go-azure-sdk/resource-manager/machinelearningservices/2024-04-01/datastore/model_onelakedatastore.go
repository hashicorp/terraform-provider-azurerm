package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Datastore = OneLakeDatastore{}

type OneLakeDatastore struct {
	Artifact                      OneLakeArtifact                `json:"artifact"`
	Endpoint                      *string                        `json:"endpoint,omitempty"`
	OneLakeWorkspaceName          string                         `json:"oneLakeWorkspaceName"`
	ServiceDataAccessAuthIdentity *ServiceDataAccessAuthIdentity `json:"serviceDataAccessAuthIdentity,omitempty"`

	// Fields inherited from Datastore
	Credentials DatastoreCredentials `json:"credentials"`
	Description *string              `json:"description,omitempty"`
	IsDefault   *bool                `json:"isDefault,omitempty"`
	Properties  *map[string]string   `json:"properties,omitempty"`
	Tags        *map[string]string   `json:"tags,omitempty"`
}

var _ json.Marshaler = OneLakeDatastore{}

func (s OneLakeDatastore) MarshalJSON() ([]byte, error) {
	type wrapper OneLakeDatastore
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OneLakeDatastore: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OneLakeDatastore: %+v", err)
	}
	decoded["datastoreType"] = "OneLake"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OneLakeDatastore: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &OneLakeDatastore{}

func (s *OneLakeDatastore) UnmarshalJSON(bytes []byte) error {
	type alias OneLakeDatastore
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into OneLakeDatastore: %+v", err)
	}

	s.Description = decoded.Description
	s.Endpoint = decoded.Endpoint
	s.IsDefault = decoded.IsDefault
	s.OneLakeWorkspaceName = decoded.OneLakeWorkspaceName
	s.Properties = decoded.Properties
	s.ServiceDataAccessAuthIdentity = decoded.ServiceDataAccessAuthIdentity
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OneLakeDatastore into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["artifact"]; ok {
		impl, err := unmarshalOneLakeArtifactImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Artifact' for 'OneLakeDatastore': %+v", err)
		}
		s.Artifact = impl
	}

	if v, ok := temp["credentials"]; ok {
		impl, err := unmarshalDatastoreCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credentials' for 'OneLakeDatastore': %+v", err)
		}
		s.Credentials = impl
	}
	return nil
}
