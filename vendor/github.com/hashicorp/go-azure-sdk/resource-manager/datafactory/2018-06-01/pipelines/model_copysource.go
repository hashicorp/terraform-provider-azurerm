package pipelines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopySource interface {
	CopySource() BaseCopySourceImpl
}

var _ CopySource = BaseCopySourceImpl{}

type BaseCopySourceImpl struct {
	DisableMetricsCollection *bool   `json:"disableMetricsCollection,omitempty"`
	MaxConcurrentConnections *int64  `json:"maxConcurrentConnections,omitempty"`
	SourceRetryCount         *int64  `json:"sourceRetryCount,omitempty"`
	SourceRetryWait          *string `json:"sourceRetryWait,omitempty"`
	Type                     string  `json:"type"`
}

func (s BaseCopySourceImpl) CopySource() BaseCopySourceImpl {
	return s
}

var _ CopySource = RawCopySourceImpl{}

// RawCopySourceImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCopySourceImpl struct {
	copySource BaseCopySourceImpl
	Type       string
	Values     map[string]interface{}
}

func (s RawCopySourceImpl) CopySource() BaseCopySourceImpl {
	return s.copySource
}

func UnmarshalCopySourceImplementation(input []byte) (CopySource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CopySource into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AmazonMWSSource") {
		var out AmazonMWSSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonMWSSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRdsForOracleSource") {
		var out AmazonRdsForOracleSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRdsForOracleSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRdsForSqlServerSource") {
		var out AmazonRdsForSqlServerSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRdsForSqlServerSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRedshiftSource") {
		var out AmazonRedshiftSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRedshiftSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AvroSource") {
		var out AvroSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AvroSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobFSSource") {
		var out AzureBlobFSSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobFSSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataExplorerSource") {
		var out AzureDataExplorerSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataExplorerSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeStoreSource") {
		var out AzureDataLakeStoreSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeStoreSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDatabricksDeltaLakeSource") {
		var out AzureDatabricksDeltaLakeSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDatabricksDeltaLakeSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMariaDBSource") {
		var out AzureMariaDBSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMariaDBSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMySqlSource") {
		var out AzureMySqlSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMySqlSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzurePostgreSqlSource") {
		var out AzurePostgreSqlSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzurePostgreSqlSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlSource") {
		var out AzureSqlSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureTableSource") {
		var out AzureTableSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureTableSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "BinarySource") {
		var out BinarySource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BinarySource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "BlobSource") {
		var out BlobSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CassandraSource") {
		var out CassandraSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CassandraSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CommonDataServiceForAppsSource") {
		var out CommonDataServiceForAppsSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CommonDataServiceForAppsSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ConcurSource") {
		var out ConcurSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConcurSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CosmosDbMongoDbApiSource") {
		var out CosmosDbMongoDbApiSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbMongoDbApiSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CosmosDbSqlApiSource") {
		var out CosmosDbSqlApiSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbSqlApiSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CouchbaseSource") {
		var out CouchbaseSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CouchbaseSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Db2Source") {
		var out Db2Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Db2Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DelimitedTextSource") {
		var out DelimitedTextSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DelimitedTextSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DocumentDbCollectionSource") {
		var out DocumentDbCollectionSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DocumentDbCollectionSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DrillSource") {
		var out DrillSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DrillSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsAXSource") {
		var out DynamicsAXSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsAXSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsCrmSource") {
		var out DynamicsCrmSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsCrmSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsSource") {
		var out DynamicsSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EloquaSource") {
		var out EloquaSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EloquaSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ExcelSource") {
		var out ExcelSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ExcelSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileSystemSource") {
		var out FileSystemSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileSystemSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleAdWordsSource") {
		var out GoogleAdWordsSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleAdWordsSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleBigQuerySource") {
		var out GoogleBigQuerySource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleBigQuerySource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleBigQueryV2Source") {
		var out GoogleBigQueryV2Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleBigQueryV2Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GreenplumSource") {
		var out GreenplumSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GreenplumSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HBaseSource") {
		var out HBaseSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HBaseSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HttpSource") {
		var out HTTPSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HTTPSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HdfsSource") {
		var out HdfsSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HdfsSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HiveSource") {
		var out HiveSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HiveSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HubspotSource") {
		var out HubspotSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HubspotSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImpalaSource") {
		var out ImpalaSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImpalaSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InformixSource") {
		var out InformixSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InformixSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "JiraSource") {
		var out JiraSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JiraSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "JsonSource") {
		var out JsonSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JsonSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LakeHouseTableSource") {
		var out LakeHouseTableSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LakeHouseTableSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MagentoSource") {
		var out MagentoSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MagentoSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MariaDBSource") {
		var out MariaDBSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MariaDBSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MarketoSource") {
		var out MarketoSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MarketoSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MicrosoftAccessSource") {
		var out MicrosoftAccessSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftAccessSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbAtlasSource") {
		var out MongoDbAtlasSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbAtlasSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbSource") {
		var out MongoDbSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbV2Source") {
		var out MongoDbV2Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbV2Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MySqlSource") {
		var out MySqlSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MySqlSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NetezzaSource") {
		var out NetezzaSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NetezzaSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ODataSource") {
		var out ODataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ODataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OdbcSource") {
		var out OdbcSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OdbcSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Office365Source") {
		var out Office365Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Office365Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleServiceCloudSource") {
		var out OracleServiceCloudSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleServiceCloudSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleSource") {
		var out OracleSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OrcSource") {
		var out OrcSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OrcSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ParquetSource") {
		var out ParquetSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ParquetSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PaypalSource") {
		var out PaypalSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PaypalSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PhoenixSource") {
		var out PhoenixSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PhoenixSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PostgreSqlSource") {
		var out PostgreSqlSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostgreSqlSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PostgreSqlV2Source") {
		var out PostgreSqlV2Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostgreSqlV2Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PrestoSource") {
		var out PrestoSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PrestoSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "QuickBooksSource") {
		var out QuickBooksSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into QuickBooksSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RelationalSource") {
		var out RelationalSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RelationalSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ResponsysSource") {
		var out ResponsysSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ResponsysSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RestSource") {
		var out RestSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RestSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceMarketingCloudSource") {
		var out SalesforceMarketingCloudSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceMarketingCloudSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceServiceCloudSource") {
		var out SalesforceServiceCloudSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceServiceCloudSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceServiceCloudV2Source") {
		var out SalesforceServiceCloudV2Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceServiceCloudV2Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceSource") {
		var out SalesforceSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceV2Source") {
		var out SalesforceV2Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceV2Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapBwSource") {
		var out SapBwSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapBwSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapCloudForCustomerSource") {
		var out SapCloudForCustomerSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapCloudForCustomerSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapEccSource") {
		var out SapEccSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapEccSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapHanaSource") {
		var out SapHanaSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapHanaSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapOdpSource") {
		var out SapOdpSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapOdpSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapOpenHubSource") {
		var out SapOpenHubSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapOpenHubSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapTableSource") {
		var out SapTableSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapTableSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceNowSource") {
		var out ServiceNowSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceNowSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceNowV2Source") {
		var out ServiceNowV2Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceNowV2Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SharePointOnlineListSource") {
		var out SharePointOnlineListSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SharePointOnlineListSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ShopifySource") {
		var out ShopifySource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ShopifySource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeSource") {
		var out SnowflakeSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeV2Source") {
		var out SnowflakeV2Source
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeV2Source: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SparkSource") {
		var out SparkSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SparkSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlDWSource") {
		var out SqlDWSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlDWSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlMISource") {
		var out SqlMISource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlMISource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlServerSource") {
		var out SqlServerSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlServerSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlSource") {
		var out SqlSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SquareSource") {
		var out SquareSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SquareSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SybaseSource") {
		var out SybaseSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SybaseSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TabularSource") {
		var out TabularSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TabularSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TeradataSource") {
		var out TeradataSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TeradataSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VerticaSource") {
		var out VerticaSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VerticaSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WarehouseSource") {
		var out WarehouseSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WarehouseSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WebSource") {
		var out WebSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "XeroSource") {
		var out XeroSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into XeroSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "XmlSource") {
		var out XmlSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into XmlSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ZohoSource") {
		var out ZohoSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ZohoSource: %+v", err)
		}
		return out, nil
	}

	var parent BaseCopySourceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCopySourceImpl: %+v", err)
	}

	return RawCopySourceImpl{
		copySource: parent,
		Type:       value,
		Values:     temp,
	}, nil

}
