package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectorServiceTypeInfoBase = AzureCosmosDBSinkConnectorServiceInfo{}

type AzureCosmosDBSinkConnectorServiceInfo struct {
	CosmosConnectionEndpoint     *string `json:"cosmosConnectionEndpoint,omitempty"`
	CosmosContainersTopicMapping *string `json:"cosmosContainersTopicMapping,omitempty"`
	CosmosDatabaseName           *string `json:"cosmosDatabaseName,omitempty"`
	CosmosIdStrategy             *string `json:"cosmosIdStrategy,omitempty"`
	CosmosMasterKey              *string `json:"cosmosMasterKey,omitempty"`

	// Fields inherited from ConnectorServiceTypeInfoBase

	ConnectorServiceType ConnectorServiceType `json:"connectorServiceType"`
}

func (s AzureCosmosDBSinkConnectorServiceInfo) ConnectorServiceTypeInfoBase() BaseConnectorServiceTypeInfoBaseImpl {
	return BaseConnectorServiceTypeInfoBaseImpl{
		ConnectorServiceType: s.ConnectorServiceType,
	}
}

var _ json.Marshaler = AzureCosmosDBSinkConnectorServiceInfo{}

func (s AzureCosmosDBSinkConnectorServiceInfo) MarshalJSON() ([]byte, error) {
	type wrapper AzureCosmosDBSinkConnectorServiceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureCosmosDBSinkConnectorServiceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureCosmosDBSinkConnectorServiceInfo: %+v", err)
	}

	decoded["connectorServiceType"] = "AzureCosmosDBSinkConnector"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureCosmosDBSinkConnectorServiceInfo: %+v", err)
	}

	return encoded, nil
}
