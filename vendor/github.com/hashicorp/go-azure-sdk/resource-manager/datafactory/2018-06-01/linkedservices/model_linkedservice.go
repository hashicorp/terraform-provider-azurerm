package linkedservices

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedService interface {
	LinkedService() BaseLinkedServiceImpl
}

var _ LinkedService = BaseLinkedServiceImpl{}

type BaseLinkedServiceImpl struct {
	Annotations *[]interface{}                     `json:"annotations,omitempty"`
	ConnectVia  *IntegrationRuntimeReference       `json:"connectVia,omitempty"`
	Description *string                            `json:"description,omitempty"`
	Parameters  *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Type        string                             `json:"type"`
	Version     *string                            `json:"version,omitempty"`
}

func (s BaseLinkedServiceImpl) LinkedService() BaseLinkedServiceImpl {
	return s
}

var _ LinkedService = RawLinkedServiceImpl{}

// RawLinkedServiceImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawLinkedServiceImpl struct {
	linkedService BaseLinkedServiceImpl
	Type          string
	Values        map[string]interface{}
}

func (s RawLinkedServiceImpl) LinkedService() BaseLinkedServiceImpl {
	return s.linkedService
}

func UnmarshalLinkedServiceImplementation(input []byte) (LinkedService, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling LinkedService into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AmazonMWS") {
		var out AmazonMWSLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonMWSLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRdsForOracle") {
		var out AmazonRdsForOracleLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRdsForOracleLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRdsForSqlServer") {
		var out AmazonRdsForSqlServerLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRdsForSqlServerLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonRedshift") {
		var out AmazonRedshiftLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonRedshiftLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonS3Compatible") {
		var out AmazonS3CompatibleLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonS3CompatibleLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AmazonS3") {
		var out AmazonS3LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AmazonS3LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AppFigures") {
		var out AppFiguresLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AppFiguresLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Asana") {
		var out AsanaLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AsanaLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBatch") {
		var out AzureBatchLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBatchLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobFS") {
		var out AzureBlobFSLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobFSLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBlobStorage") {
		var out AzureBlobStorageLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobStorageLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataExplorer") {
		var out AzureDataExplorerLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataExplorerLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeAnalytics") {
		var out AzureDataLakeAnalyticsLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeAnalyticsLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeStore") {
		var out AzureDataLakeStoreLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeStoreLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDatabricksDeltaLake") {
		var out AzureDatabricksDeltaLakeLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDatabricksDeltaLakeLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDatabricks") {
		var out AzureDatabricksLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDatabricksLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureFileStorage") {
		var out AzureFileStorageLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileStorageLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureFunction") {
		var out AzureFunctionLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFunctionLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureKeyVault") {
		var out AzureKeyVaultLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureKeyVaultLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureML") {
		var out AzureMLLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMLLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMLService") {
		var out AzureMLServiceLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMLServiceLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMariaDB") {
		var out AzureMariaDBLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMariaDBLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureMySql") {
		var out AzureMySqlLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMySqlLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzurePostgreSql") {
		var out AzurePostgreSqlLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzurePostgreSqlLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSearch") {
		var out AzureSearchLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSearchLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlDW") {
		var out AzureSqlDWLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlDWLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlDatabase") {
		var out AzureSqlDatabaseLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlDatabaseLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSqlMI") {
		var out AzureSqlMILinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSqlMILinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureStorage") {
		var out AzureStorageLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureStorageLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureSynapseArtifacts") {
		var out AzureSynapseArtifactsLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureSynapseArtifactsLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureTableStorage") {
		var out AzureTableStorageLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureTableStorageLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Cassandra") {
		var out CassandraLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CassandraLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CommonDataServiceForApps") {
		var out CommonDataServiceForAppsLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CommonDataServiceForAppsLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Concur") {
		var out ConcurLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConcurLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CosmosDb") {
		var out CosmosDbLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CosmosDbMongoDbApi") {
		var out CosmosDbMongoDbApiLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbMongoDbApiLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Couchbase") {
		var out CouchbaseLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CouchbaseLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CustomDataSource") {
		var out CustomDataSourceLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomDataSourceLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Dataworld") {
		var out DataworldLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataworldLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Db2") {
		var out Db2LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Db2LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Drill") {
		var out DrillLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DrillLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsAX") {
		var out DynamicsAXLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsAXLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DynamicsCrm") {
		var out DynamicsCrmLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsCrmLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Dynamics") {
		var out DynamicsLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicsLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Eloqua") {
		var out EloquaLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EloquaLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FileServer") {
		var out FileServerLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FileServerLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FtpServer") {
		var out FtpServerLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FtpServerLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleAdWords") {
		var out GoogleAdWordsLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleAdWordsLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleBigQuery") {
		var out GoogleBigQueryLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleBigQueryLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleBigQueryV2") {
		var out GoogleBigQueryV2LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleBigQueryV2LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleCloudStorage") {
		var out GoogleCloudStorageLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleCloudStorageLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GoogleSheets") {
		var out GoogleSheetsLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GoogleSheetsLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Greenplum") {
		var out GreenplumLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GreenplumLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HBase") {
		var out HBaseLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HBaseLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HDInsight") {
		var out HDInsightLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HDInsightLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HDInsightOnDemand") {
		var out HDInsightOnDemandLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HDInsightOnDemandLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HttpServer") {
		var out HTTPLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HTTPLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Hdfs") {
		var out HdfsLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HdfsLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Hive") {
		var out HiveLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HiveLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Hubspot") {
		var out HubspotLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HubspotLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Impala") {
		var out ImpalaLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImpalaLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Informix") {
		var out InformixLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InformixLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Jira") {
		var out JiraLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JiraLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Lakehouse") {
		var out LakeHouseLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LakeHouseLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Magento") {
		var out MagentoLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MagentoLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MariaDB") {
		var out MariaDBLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MariaDBLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Marketo") {
		var out MarketoLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MarketoLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MicrosoftAccess") {
		var out MicrosoftAccessLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MicrosoftAccessLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbAtlas") {
		var out MongoDbAtlasLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbAtlasLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDb") {
		var out MongoDbLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MongoDbV2") {
		var out MongoDbV2LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbV2LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MySql") {
		var out MySqlLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MySqlLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Netezza") {
		var out NetezzaLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NetezzaLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OData") {
		var out ODataLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ODataLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Odbc") {
		var out OdbcLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OdbcLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Office365") {
		var out Office365LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Office365LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleCloudStorage") {
		var out OracleCloudStorageLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleCloudStorageLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Oracle") {
		var out OracleLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OracleServiceCloud") {
		var out OracleServiceCloudLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OracleServiceCloudLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Paypal") {
		var out PaypalLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PaypalLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Phoenix") {
		var out PhoenixLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PhoenixLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PostgreSql") {
		var out PostgreSqlLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostgreSqlLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PostgreSqlV2") {
		var out PostgreSqlV2LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostgreSqlV2LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Presto") {
		var out PrestoLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PrestoLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "QuickBooks") {
		var out QuickBooksLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into QuickBooksLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Quickbase") {
		var out QuickbaseLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into QuickbaseLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Responsys") {
		var out ResponsysLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ResponsysLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RestService") {
		var out RestServiceLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RestServiceLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Salesforce") {
		var out SalesforceLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceMarketingCloud") {
		var out SalesforceMarketingCloudLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceMarketingCloudLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceServiceCloud") {
		var out SalesforceServiceCloudLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceServiceCloudLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceServiceCloudV2") {
		var out SalesforceServiceCloudV2LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceServiceCloudV2LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SalesforceV2") {
		var out SalesforceV2LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SalesforceV2LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapBW") {
		var out SapBWLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapBWLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapCloudForCustomer") {
		var out SapCloudForCustomerLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapCloudForCustomerLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapEcc") {
		var out SapEccLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapEccLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapHana") {
		var out SapHanaLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapHanaLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapOdp") {
		var out SapOdpLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapOdpLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapOpenHub") {
		var out SapOpenHubLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapOpenHubLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SapTable") {
		var out SapTableLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SapTableLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceNow") {
		var out ServiceNowLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceNowLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceNowV2") {
		var out ServiceNowV2LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceNowV2LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Sftp") {
		var out SftpServerLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SftpServerLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SharePointOnlineList") {
		var out SharePointOnlineListLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SharePointOnlineListLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Shopify") {
		var out ShopifyLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ShopifyLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Smartsheet") {
		var out SmartsheetLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SmartsheetLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Snowflake") {
		var out SnowflakeLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SnowflakeV2") {
		var out SnowflakeV2LinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SnowflakeV2LinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Spark") {
		var out SparkLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SparkLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlServer") {
		var out SqlServerLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlServerLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Square") {
		var out SquareLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SquareLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Sybase") {
		var out SybaseLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SybaseLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TeamDesk") {
		var out TeamDeskLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TeamDeskLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Teradata") {
		var out TeradataLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TeradataLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Twilio") {
		var out TwilioLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TwilioLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Vertica") {
		var out VerticaLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VerticaLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Warehouse") {
		var out WarehouseLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WarehouseLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Web") {
		var out WebLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Xero") {
		var out XeroLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into XeroLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Zendesk") {
		var out ZendeskLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ZendeskLinkedService: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Zoho") {
		var out ZohoLinkedService
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ZohoLinkedService: %+v", err)
		}
		return out, nil
	}

	var parent BaseLinkedServiceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseLinkedServiceImpl: %+v", err)
	}

	return RawLinkedServiceImpl{
		linkedService: parent,
		Type:          value,
		Values:        temp,
	}, nil

}
