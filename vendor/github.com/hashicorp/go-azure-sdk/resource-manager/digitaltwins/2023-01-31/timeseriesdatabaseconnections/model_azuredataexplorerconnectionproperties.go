package timeseriesdatabaseconnections

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TimeSeriesDatabaseConnectionProperties = AzureDataExplorerConnectionProperties{}

type AzureDataExplorerConnectionProperties struct {
	AdxDatabaseName                         string  `json:"adxDatabaseName"`
	AdxEndpointUri                          string  `json:"adxEndpointUri"`
	AdxRelationshipLifecycleEventsTableName *string `json:"adxRelationshipLifecycleEventsTableName,omitempty"`
	AdxResourceId                           string  `json:"adxResourceId"`
	AdxTableName                            *string `json:"adxTableName,omitempty"`
	AdxTwinLifecycleEventsTableName         *string `json:"adxTwinLifecycleEventsTableName,omitempty"`
	EventHubConsumerGroup                   *string `json:"eventHubConsumerGroup,omitempty"`
	EventHubEndpointUri                     string  `json:"eventHubEndpointUri"`
	EventHubEntityPath                      string  `json:"eventHubEntityPath"`
	EventHubNamespaceResourceId             string  `json:"eventHubNamespaceResourceId"`
	RecordPropertyAndItemRemovals           *bool   `json:"recordPropertyAndItemRemovals,omitempty"`

	// Fields inherited from TimeSeriesDatabaseConnectionProperties

	ConnectionType    ConnectionType                     `json:"connectionType"`
	Identity          *ManagedIdentityReference          `json:"identity,omitempty"`
	ProvisioningState *TimeSeriesDatabaseConnectionState `json:"provisioningState,omitempty"`
}

func (s AzureDataExplorerConnectionProperties) TimeSeriesDatabaseConnectionProperties() BaseTimeSeriesDatabaseConnectionPropertiesImpl {
	return BaseTimeSeriesDatabaseConnectionPropertiesImpl{
		ConnectionType:    s.ConnectionType,
		Identity:          s.Identity,
		ProvisioningState: s.ProvisioningState,
	}
}

var _ json.Marshaler = AzureDataExplorerConnectionProperties{}

func (s AzureDataExplorerConnectionProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureDataExplorerConnectionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDataExplorerConnectionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDataExplorerConnectionProperties: %+v", err)
	}

	decoded["connectionType"] = "AzureDataExplorer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDataExplorerConnectionProperties: %+v", err)
	}

	return encoded, nil
}
