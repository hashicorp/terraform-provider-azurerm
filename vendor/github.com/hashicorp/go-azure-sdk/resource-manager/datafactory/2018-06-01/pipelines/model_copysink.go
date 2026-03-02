package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopySink interface {
	CopySink() BaseCopySinkImpl
}

var _ CopySink = BaseCopySinkImpl{}

type BaseCopySinkImpl struct {
	DisableMetricsCollection *bool        `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64       `json:"maxConcurrentConnections,omitempty"`
	SinkRetryCount           *int64       `json:"sinkRetryCount,omitempty"`
	SinkRetryWait            *interface{} `json:"sinkRetryWait,omitempty"`
	Type                     string       `json:"type"`
	WriteBatchSize           *int64       `json:"writeBatchSize,omitempty"`
	WriteBatchTimeout        *interface{} `json:"writeBatchTimeout,omitempty"`
}

func (s BaseCopySinkImpl) CopySink() BaseCopySinkImpl {
	return s
}

var _ CopySink = RawCopySinkImpl{}

// RawCopySinkImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCopySinkImpl struct {
	copySink BaseCopySinkImpl
	Type     string
	Values   map[string]interface{}
}

func (s RawCopySinkImpl) CopySink() BaseCopySinkImpl {
	return s.copySink
}

func UnmarshalCopySinkImplementation(input []byte) (CopySink, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CopySink into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AvroSink") {
		var out AvroSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AvroSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobFSSink") {
		var out AzureBlobFSSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobFSSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataExplorerSink") {
		var out AzureDataExplorerSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataExplorerSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeStoreSink") {
		var out AzureDataLakeStoreSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeStoreSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDatabricksDeltaLakeSink") {
		var out AzureDatabricksDeltaLakeSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDatabricksDeltaLakeSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMySqlSink") {
		var out AzureMySqlSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMySqlSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzurePostgreSqlSink") {
		var out AzurePostgreSqlSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzurePostgreSqlSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureQueueSink") {
		var out AzureQueueSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureQueueSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSearchIndexSink") {
		var out AzureSearchIndexSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSearchIndexSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlSink") {
		var out AzureSqlSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureTableSink") {
		var out AzureTableSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureTableSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "BinarySink") {
		var out BinarySink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BinarySink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "BlobSink") {
		var out BlobSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CommonDataServiceForAppsSink") {
		var out CommonDataServiceForAppsSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CommonDataServiceForAppsSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CosmosDbMongoDbApiSink") {
		var out CosmosDbMongoDbApiSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbMongoDbApiSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CosmosDbSqlApiSink") {
		var out CosmosDbSqlApiSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbSqlApiSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DelimitedTextSink") {
		var out DelimitedTextSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DelimitedTextSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DocumentDbCollectionSink") {
		var out DocumentDbCollectionSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DocumentDbCollectionSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsCrmSink") {
		var out DynamicsCrmSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsCrmSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsSink") {
		var out DynamicsSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileSystemSink") {
		var out FileSystemSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileSystemSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IcebergSink") {
		var out IcebergSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IcebergSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InformixSink") {
		var out InformixSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InformixSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "JsonSink") {
		var out JsonSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JsonSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LakeHouseTableSink") {
		var out LakeHouseTableSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LakeHouseTableSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MicrosoftAccessSink") {
		var out MicrosoftAccessSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftAccessSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbAtlasSink") {
		var out MongoDbAtlasSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbAtlasSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbV2Sink") {
		var out MongoDbV2Sink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbV2Sink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OdbcSink") {
		var out OdbcSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OdbcSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleSink") {
		var out OracleSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OrcSink") {
		var out OrcSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OrcSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ParquetSink") {
		var out ParquetSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ParquetSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RestSink") {
		var out RestSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RestSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceServiceCloudSink") {
		var out SalesforceServiceCloudSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceServiceCloudSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceServiceCloudV2Sink") {
		var out SalesforceServiceCloudV2Sink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceServiceCloudV2Sink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceSink") {
		var out SalesforceSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceV2Sink") {
		var out SalesforceV2Sink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceV2Sink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapCloudForCustomerSink") {
		var out SapCloudForCustomerSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapCloudForCustomerSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeSink") {
		var out SnowflakeSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeV2Sink") {
		var out SnowflakeV2Sink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeV2Sink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlDWSink") {
		var out SqlDWSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlDWSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlMISink") {
		var out SqlMISink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlMISink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlServerSink") {
		var out SqlServerSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlServerSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlSink") {
		var out SqlSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TeradataSink") {
		var out TeradataSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TeradataSink: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WarehouseSink") {
		var out WarehouseSink
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WarehouseSink: %+v", err)
		}
		return out, nil
	}

	var parent BaseCopySinkImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCopySinkImpl: %+v", err)
	}

	return RawCopySinkImpl{
		copySink: parent,
		Type:     value,
		Values:   temp,
	}, nil

}
