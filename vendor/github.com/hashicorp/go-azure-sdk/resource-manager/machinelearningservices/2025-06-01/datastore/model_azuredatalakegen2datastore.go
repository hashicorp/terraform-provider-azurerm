package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Datastore = AzureDataLakeGen2Datastore{}

type AzureDataLakeGen2Datastore struct {
	AccountName                   string                         `json:"accountName"`
	Endpoint                      *string                        `json:"endpoint,omitempty"`
	Filesystem                    string                         `json:"filesystem"`
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

func (s AzureDataLakeGen2Datastore) Datastore() BaseDatastoreImpl {
	return BaseDatastoreImpl{
		Credentials:   s.Credentials,
		DatastoreType: s.DatastoreType,
		Description:   s.Description,
		IsDefault:     s.IsDefault,
		Properties:    s.Properties,
		Tags:          s.Tags,
	}
}

var _ json.Marshaler = AzureDataLakeGen2Datastore{}

func (s AzureDataLakeGen2Datastore) MarshalJSON() ([]byte, error) {
	type wrapper AzureDataLakeGen2Datastore
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDataLakeGen2Datastore: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDataLakeGen2Datastore: %+v", err)
	}

	decoded["datastoreType"] = "AzureDataLakeGen2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDataLakeGen2Datastore: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureDataLakeGen2Datastore{}

func (s *AzureDataLakeGen2Datastore) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccountName                   string                         `json:"accountName"`
		Endpoint                      *string                        `json:"endpoint,omitempty"`
		Filesystem                    string                         `json:"filesystem"`
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
	s.Filesystem = decoded.Filesystem
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
		return fmt.Errorf("unmarshaling AzureDataLakeGen2Datastore into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["credentials"]; ok {
		impl, err := UnmarshalDatastoreCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credentials' for 'AzureDataLakeGen2Datastore': %+v", err)
		}
		s.Credentials = impl
	}

	return nil
}
