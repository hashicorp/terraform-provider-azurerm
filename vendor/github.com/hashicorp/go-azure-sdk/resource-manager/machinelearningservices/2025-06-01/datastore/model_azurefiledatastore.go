package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Datastore = AzureFileDatastore{}

type AzureFileDatastore struct {
	AccountName                   string                         `json:"accountName"`
	Endpoint                      *string                        `json:"endpoint,omitempty"`
	FileShareName                 string                         `json:"fileShareName"`
	Protocol                      *string                        `json:"protocol,omitempty"`
	ResourceGroup                 *string                        `json:"resourceGroup,omitempty"`
	ServiceDataAccessAuthIdentity *ServiceDataAccessAuthIdentity `json:"serviceDataAccessAuthIdentity,omitempty"`
	SubscriptionId                *string                        `json:"subscriptionId,omitempty"`

	// Fields inherited from Datastore

	Credentials   DatastoreCredentials `json:"credentials"`
	DatastoreType DatastoreType        `json:"datastoreType"`
	Description   *string              `json:"description,omitempty"`
	IsDefault     *bool                `json:"isDefault,omitempty"`
	Properties    *map[string]string   `json:"properties,omitempty"`
	Tags          *map[string]string   `json:"tags,omitempty"`
}

func (s AzureFileDatastore) Datastore() BaseDatastoreImpl {
	return BaseDatastoreImpl{
		Credentials:   s.Credentials,
		DatastoreType: s.DatastoreType,
		Description:   s.Description,
		IsDefault:     s.IsDefault,
		Properties:    s.Properties,
		Tags:          s.Tags,
	}
}

var _ json.Marshaler = AzureFileDatastore{}

func (s AzureFileDatastore) MarshalJSON() ([]byte, error) {
	type wrapper AzureFileDatastore
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFileDatastore: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFileDatastore: %+v", err)
	}

	decoded["datastoreType"] = "AzureFile"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFileDatastore: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureFileDatastore{}

func (s *AzureFileDatastore) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountName                   string                         `json:"accountName"`
		Endpoint                      *string                        `json:"endpoint,omitempty"`
		FileShareName                 string                         `json:"fileShareName"`
		Protocol                      *string                        `json:"protocol,omitempty"`
		ResourceGroup                 *string                        `json:"resourceGroup,omitempty"`
		ServiceDataAccessAuthIdentity *ServiceDataAccessAuthIdentity `json:"serviceDataAccessAuthIdentity,omitempty"`
		SubscriptionId                *string                        `json:"subscriptionId,omitempty"`
		DatastoreType                 DatastoreType                  `json:"datastoreType"`
		Description                   *string                        `json:"description,omitempty"`
		IsDefault                     *bool                          `json:"isDefault,omitempty"`
		Properties                    *map[string]string             `json:"properties,omitempty"`
		Tags                          *map[string]string             `json:"tags,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccountName = decoded.AccountName
	s.Endpoint = decoded.Endpoint
	s.FileShareName = decoded.FileShareName
	s.Protocol = decoded.Protocol
	s.ResourceGroup = decoded.ResourceGroup
	s.ServiceDataAccessAuthIdentity = decoded.ServiceDataAccessAuthIdentity
	s.SubscriptionId = decoded.SubscriptionId
	s.DatastoreType = decoded.DatastoreType
	s.Description = decoded.Description
	s.IsDefault = decoded.IsDefault
	s.Properties = decoded.Properties
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureFileDatastore into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["credentials"]; ok {
		impl, err := UnmarshalDatastoreCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credentials' for 'AzureFileDatastore': %+v", err)
		}
		s.Credentials = impl
	}

	return nil
}
