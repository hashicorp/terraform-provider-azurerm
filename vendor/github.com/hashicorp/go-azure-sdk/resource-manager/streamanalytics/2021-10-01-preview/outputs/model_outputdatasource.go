package outputs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutputDataSource interface {
	OutputDataSource() BaseOutputDataSourceImpl
}

var _ OutputDataSource = BaseOutputDataSourceImpl{}

type BaseOutputDataSourceImpl struct {
	Type string `json:"type"`
}

func (s BaseOutputDataSourceImpl) OutputDataSource() BaseOutputDataSourceImpl {
	return s
}

var _ OutputDataSource = RawOutputDataSourceImpl{}

// RawOutputDataSourceImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawOutputDataSourceImpl struct {
	outputDataSource BaseOutputDataSourceImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawOutputDataSourceImpl) OutputDataSource() BaseOutputDataSourceImpl {
	return s.outputDataSource
}

func UnmarshalOutputDataSourceImplementation(input []byte) (OutputDataSource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling OutputDataSource into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Microsoft.Kusto/clusters/databases") {
		var out AzureDataExplorerOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataExplorerOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.DataLake/Accounts") {
		var out AzureDataLakeStoreOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeStoreOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.AzureFunction") {
		var out AzureFunctionOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFunctionOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Sql/Server/Database") {
		var out AzureSqlDatabaseOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlDatabaseOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Sql/Server/DataWarehouse") {
		var out AzureSynapseOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSynapseOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Storage/Table") {
		var out AzureTableOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureTableOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Storage/Blob") {
		var out BlobOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.Storage/DocumentDB") {
		var out DocumentDbOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DocumentDbOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.ServiceBus/EventHub") {
		var out EventHubOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.EventHub/EventHub") {
		var out EventHubV2OutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubV2OutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GatewayMessageBus") {
		var out GatewayMessageBusOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GatewayMessageBusOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.DBForPostgreSQL/servers/databases") {
		var out PostgreSQLOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostgreSQLOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PowerBI") {
		var out PowerBIOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PowerBIOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Raw") {
		var out RawOutputDatasource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RawOutputDatasource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.ServiceBus/Queue") {
		var out ServiceBusQueueOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceBusQueueOutputDataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.ServiceBus/Topic") {
		var out ServiceBusTopicOutputDataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceBusTopicOutputDataSource: %+v", err)
		}
		return out, nil
	}

	var parent BaseOutputDataSourceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseOutputDataSourceImpl: %+v", err)
	}

	return RawOutputDataSourceImpl{
		outputDataSource: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
