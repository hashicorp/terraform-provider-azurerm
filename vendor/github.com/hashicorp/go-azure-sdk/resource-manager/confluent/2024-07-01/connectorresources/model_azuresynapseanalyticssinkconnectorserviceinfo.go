package connectorresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectorServiceTypeInfoBase = AzureSynapseAnalyticsSinkConnectorServiceInfo{}

type AzureSynapseAnalyticsSinkConnectorServiceInfo struct {
	SynapseSqlDatabaseName *string `json:"synapseSqlDatabaseName,omitempty"`
	SynapseSqlPassword     *string `json:"synapseSqlPassword,omitempty"`
	SynapseSqlServerName   *string `json:"synapseSqlServerName,omitempty"`
	SynapseSqlUser         *string `json:"synapseSqlUser,omitempty"`

	// Fields inherited from ConnectorServiceTypeInfoBase

	ConnectorServiceType ConnectorServiceType `json:"connectorServiceType"`
}

func (s AzureSynapseAnalyticsSinkConnectorServiceInfo) ConnectorServiceTypeInfoBase() BaseConnectorServiceTypeInfoBaseImpl {
	return BaseConnectorServiceTypeInfoBaseImpl{
		ConnectorServiceType: s.ConnectorServiceType,
	}
}

var _ json.Marshaler = AzureSynapseAnalyticsSinkConnectorServiceInfo{}

func (s AzureSynapseAnalyticsSinkConnectorServiceInfo) MarshalJSON() ([]byte, error) {
	type wrapper AzureSynapseAnalyticsSinkConnectorServiceInfo
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureSynapseAnalyticsSinkConnectorServiceInfo: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureSynapseAnalyticsSinkConnectorServiceInfo: %+v", err)
	}

	decoded["connectorServiceType"] = "AzureSynapseAnalyticsSinkConnector"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureSynapseAnalyticsSinkConnectorServiceInfo: %+v", err)
	}

	return encoded, nil
}
