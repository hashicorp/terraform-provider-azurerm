package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Datastore = AzureDataLakeGen1Datastore{}

type AzureDataLakeGen1Datastore struct {
	ResourceGroup                 *string                        `json:"resourceGroup,omitempty"`
	ServiceDataAccessAuthIdentity *ServiceDataAccessAuthIdentity `json:"serviceDataAccessAuthIdentity,omitempty"`
	StoreName                     string                         `json:"storeName"`
	SubscriptionId                *string                        `json:"subscriptionId,omitempty"`

	// Fields inherited from Datastore

	Credentials   DatastoreCredentials `json:"credentials"`
	DatastoreType DatastoreType        `json:"datastoreType"`
	Description   *string              `json:"description,omitempty"`
	IsDefault     *bool                `json:"isDefault,omitempty"`
	Properties    *map[string]string   `json:"properties,omitempty"`
	Tags          *map[string]string   `json:"tags,omitempty"`
}

func (s AzureDataLakeGen1Datastore) Datastore() BaseDatastoreImpl {
	return BaseDatastoreImpl{
		Credentials:   s.Credentials,
		DatastoreType: s.DatastoreType,
		Description:   s.Description,
		IsDefault:     s.IsDefault,
		Properties:    s.Properties,
		Tags:          s.Tags,
	}
}

var _ json.Marshaler = AzureDataLakeGen1Datastore{}

func (s AzureDataLakeGen1Datastore) MarshalJSON() ([]byte, error) {
	type wrapper AzureDataLakeGen1Datastore
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDataLakeGen1Datastore: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDataLakeGen1Datastore: %+v", err)
	}

	decoded["datastoreType"] = "AzureDataLakeGen1"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDataLakeGen1Datastore: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureDataLakeGen1Datastore{}

func (s *AzureDataLakeGen1Datastore) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ResourceGroup                 *string                        `json:"resourceGroup,omitempty"`
		ServiceDataAccessAuthIdentity *ServiceDataAccessAuthIdentity `json:"serviceDataAccessAuthIdentity,omitempty"`
		StoreName                     string                         `json:"storeName"`
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

	s.ResourceGroup = decoded.ResourceGroup
	s.ServiceDataAccessAuthIdentity = decoded.ServiceDataAccessAuthIdentity
	s.StoreName = decoded.StoreName
	s.SubscriptionId = decoded.SubscriptionId
	s.DatastoreType = decoded.DatastoreType
	s.Description = decoded.Description
	s.IsDefault = decoded.IsDefault
	s.Properties = decoded.Properties
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureDataLakeGen1Datastore into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["credentials"]; ok {
		impl, err := UnmarshalDatastoreCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credentials' for 'AzureDataLakeGen1Datastore': %+v", err)
		}
		s.Credentials = impl
	}

	return nil
}
