package dataconnectors

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnector = AwsS3DataConnector{}

type AwsS3DataConnector struct {
	Properties *AwsS3DataConnectorProperties `json:"properties,omitempty"`

	// Fields inherited from DataConnector

	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Kind       DataConnectorKind      `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s AwsS3DataConnector) DataConnector() BaseDataConnectorImpl {
	return BaseDataConnectorImpl{
		Etag:       s.Etag,
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = AwsS3DataConnector{}

func (s AwsS3DataConnector) MarshalJSON() ([]byte, error) {
	type wrapper AwsS3DataConnector
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AwsS3DataConnector: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AwsS3DataConnector: %+v", err)
	}

	decoded["kind"] = "AmazonWebServicesS3"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AwsS3DataConnector: %+v", err)
	}

	return encoded, nil
}
