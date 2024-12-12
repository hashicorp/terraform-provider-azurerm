package datasets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Dataset interface {
	Dataset() BaseDatasetImpl
}

var _ Dataset = BaseDatasetImpl{}

type BaseDatasetImpl struct {
	Annotations       *[]interface{}                     `json:"annotations,omitempty"`
	Description       *string                            `json:"description,omitempty"`
	Folder            *DatasetFolder                     `json:"folder,omitempty"`
	LinkedServiceName LinkedServiceReference             `json:"linkedServiceName"`
	Parameters        *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Schema            *interface{}                       `json:"schema,omitempty"`
	Structure         *interface{}                       `json:"structure,omitempty"`
	Type              string                             `json:"type"`
}

func (s BaseDatasetImpl) Dataset() BaseDatasetImpl {
	return s
}

var _ Dataset = RawDatasetImpl{}

// RawDatasetImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDatasetImpl struct {
	dataset BaseDatasetImpl
	Type    string
	Values  map[string]interface{}
}

func (s RawDatasetImpl) Dataset() BaseDatasetImpl {
	return s.dataset
}

func UnmarshalDatasetImplementation(input []byte) (Dataset, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Dataset into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AmazonMWSObject") {
		var out AmazonMWSObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonMWSObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRdsForOracleTable") {
		var out AmazonRdsForOracleTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRdsForOracleTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRdsForSqlServerTable") {
		var out AmazonRdsForSqlServerTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRdsForSqlServerTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRedshiftTable") {
		var out AmazonRedshiftTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRedshiftTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonS3Object") {
		var out AmazonS3Dataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonS3Dataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Avro") {
		var out AvroDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AvroDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlob") {
		var out AzureBlobDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobFSFile") {
		var out AzureBlobFSDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobFSDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataExplorerTable") {
		var out AzureDataExplorerTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataExplorerTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeStoreFile") {
		var out AzureDataLakeStoreDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeStoreDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDatabricksDeltaLakeDataset") {
		var out AzureDatabricksDeltaLakeDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDatabricksDeltaLakeDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMariaDBTable") {
		var out AzureMariaDBTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMariaDBTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMySqlTable") {
		var out AzureMySqlTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMySqlTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzurePostgreSqlTable") {
		var out AzurePostgreSqlTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzurePostgreSqlTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSearchIndex") {
		var out AzureSearchIndexDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSearchIndexDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlDWTable") {
		var out AzureSqlDWTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlDWTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlMITable") {
		var out AzureSqlMITableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlMITableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlTable") {
		var out AzureSqlTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureTable") {
		var out AzureTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Binary") {
		var out BinaryDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BinaryDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CassandraTable") {
		var out CassandraTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CassandraTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CommonDataServiceForAppsEntity") {
		var out CommonDataServiceForAppsEntityDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CommonDataServiceForAppsEntityDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConcurObject") {
		var out ConcurObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConcurObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CosmosDbMongoDbApiCollection") {
		var out CosmosDbMongoDbApiCollectionDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbMongoDbApiCollectionDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CosmosDbSqlApiCollection") {
		var out CosmosDbSqlApiCollectionDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbSqlApiCollectionDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CouchbaseTable") {
		var out CouchbaseTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CouchbaseTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CustomDataset") {
		var out CustomDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Db2Table") {
		var out Db2TableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Db2TableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DelimitedText") {
		var out DelimitedTextDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DelimitedTextDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DocumentDbCollection") {
		var out DocumentDbCollectionDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DocumentDbCollectionDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DrillTable") {
		var out DrillTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DrillTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsAXResource") {
		var out DynamicsAXResourceDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsAXResourceDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsCrmEntity") {
		var out DynamicsCrmEntityDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsCrmEntityDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsEntity") {
		var out DynamicsEntityDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsEntityDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EloquaObject") {
		var out EloquaObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EloquaObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Excel") {
		var out ExcelDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExcelDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileShare") {
		var out FileShareDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileShareDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleAdWordsObject") {
		var out GoogleAdWordsObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleAdWordsObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleBigQueryObject") {
		var out GoogleBigQueryObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleBigQueryObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleBigQueryV2Object") {
		var out GoogleBigQueryV2ObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleBigQueryV2ObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GreenplumTable") {
		var out GreenplumTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GreenplumTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HBaseObject") {
		var out HBaseObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HBaseObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HttpFile") {
		var out HTTPDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HTTPDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HiveObject") {
		var out HiveObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HiveObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HubspotObject") {
		var out HubspotObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HubspotObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Iceberg") {
		var out IcebergDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IcebergDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImpalaObject") {
		var out ImpalaObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImpalaObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InformixTable") {
		var out InformixTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InformixTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "JiraObject") {
		var out JiraObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JiraObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Json") {
		var out JsonDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JsonDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LakeHouseTable") {
		var out LakeHouseTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LakeHouseTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MagentoObject") {
		var out MagentoObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MagentoObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MariaDBTable") {
		var out MariaDBTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MariaDBTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MarketoObject") {
		var out MarketoObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MarketoObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MicrosoftAccessTable") {
		var out MicrosoftAccessTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftAccessTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbAtlasCollection") {
		var out MongoDbAtlasCollectionDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbAtlasCollectionDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbCollection") {
		var out MongoDbCollectionDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbCollectionDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbV2Collection") {
		var out MongoDbV2CollectionDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbV2CollectionDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MySqlTable") {
		var out MySqlTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MySqlTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NetezzaTable") {
		var out NetezzaTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NetezzaTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ODataResource") {
		var out ODataResourceDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ODataResourceDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OdbcTable") {
		var out OdbcTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OdbcTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Office365Table") {
		var out Office365Dataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Office365Dataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleServiceCloudObject") {
		var out OracleServiceCloudObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleServiceCloudObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleTable") {
		var out OracleTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Orc") {
		var out OrcDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OrcDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Parquet") {
		var out ParquetDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ParquetDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PaypalObject") {
		var out PaypalObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PaypalObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PhoenixObject") {
		var out PhoenixObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PhoenixObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PostgreSqlTable") {
		var out PostgreSqlTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostgreSqlTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PostgreSqlV2Table") {
		var out PostgreSqlV2TableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostgreSqlV2TableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PrestoObject") {
		var out PrestoObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PrestoObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "QuickBooksObject") {
		var out QuickBooksObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into QuickBooksObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RelationalTable") {
		var out RelationalTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RelationalTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ResponsysObject") {
		var out ResponsysObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ResponsysObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RestResource") {
		var out RestResourceDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RestResourceDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceMarketingCloudObject") {
		var out SalesforceMarketingCloudObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceMarketingCloudObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceObject") {
		var out SalesforceObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceServiceCloudObject") {
		var out SalesforceServiceCloudObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceServiceCloudObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceServiceCloudV2Object") {
		var out SalesforceServiceCloudV2ObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceServiceCloudV2ObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceV2Object") {
		var out SalesforceV2ObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceV2ObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapBwCube") {
		var out SapBwCubeDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapBwCubeDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapCloudForCustomerResource") {
		var out SapCloudForCustomerResourceDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapCloudForCustomerResourceDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapEccResource") {
		var out SapEccResourceDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapEccResourceDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapHanaTable") {
		var out SapHanaTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapHanaTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapOdpResource") {
		var out SapOdpResourceDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapOdpResourceDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapOpenHubTable") {
		var out SapOpenHubTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapOpenHubTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapTableResource") {
		var out SapTableResourceDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapTableResourceDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceNowObject") {
		var out ServiceNowObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceNowObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceNowV2Object") {
		var out ServiceNowV2ObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceNowV2ObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SharePointOnlineListResource") {
		var out SharePointOnlineListResourceDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SharePointOnlineListResourceDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ShopifyObject") {
		var out ShopifyObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ShopifyObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeTable") {
		var out SnowflakeDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeV2Table") {
		var out SnowflakeV2Dataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeV2Dataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SparkObject") {
		var out SparkObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SparkObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlServerTable") {
		var out SqlServerTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlServerTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SquareObject") {
		var out SquareObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SquareObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SybaseTable") {
		var out SybaseTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SybaseTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TeradataTable") {
		var out TeradataTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TeradataTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VerticaTable") {
		var out VerticaTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VerticaTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WarehouseTable") {
		var out WarehouseTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WarehouseTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WebTable") {
		var out WebTableDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebTableDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "XeroObject") {
		var out XeroObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into XeroObjectDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Xml") {
		var out XmlDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into XmlDataset: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ZohoObject") {
		var out ZohoObjectDataset
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ZohoObjectDataset: %+v", err)
		}
		return out, nil
	}

	var parent BaseDatasetImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDatasetImpl: %+v", err)
	}

	return RawDatasetImpl{
		dataset: parent,
		Type:    value,
		Values:  temp,
	}, nil

}
