package accountconnectionresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionAuthType string

const (
	ConnectionAuthTypeAAD                    ConnectionAuthType = "AAD"
	ConnectionAuthTypeAccessKey              ConnectionAuthType = "AccessKey"
	ConnectionAuthTypeAccountKey             ConnectionAuthType = "AccountKey"
	ConnectionAuthTypeAccountManagedIdentity ConnectionAuthType = "AccountManagedIdentity"
	ConnectionAuthTypeAgentUserImpersonation ConnectionAuthType = "AgentUserImpersonation"
	ConnectionAuthTypeAgenticIdentityToken   ConnectionAuthType = "AgenticIdentityToken"
	ConnectionAuthTypeAgenticUser            ConnectionAuthType = "AgenticUser"
	ConnectionAuthTypeApiKey                 ConnectionAuthType = "ApiKey"
	ConnectionAuthTypeCustomKeys             ConnectionAuthType = "CustomKeys"
	ConnectionAuthTypeDelegatedSAS           ConnectionAuthType = "DelegatedSAS"
	ConnectionAuthTypeManagedIdentity        ConnectionAuthType = "ManagedIdentity"
	ConnectionAuthTypeNone                   ConnectionAuthType = "None"
	ConnectionAuthTypeOAuthTwo               ConnectionAuthType = "OAuth2"
	ConnectionAuthTypePAT                    ConnectionAuthType = "PAT"
	ConnectionAuthTypeProjectManagedIdentity ConnectionAuthType = "ProjectManagedIdentity"
	ConnectionAuthTypeSAS                    ConnectionAuthType = "SAS"
	ConnectionAuthTypeServicePrincipal       ConnectionAuthType = "ServicePrincipal"
	ConnectionAuthTypeUserEntraToken         ConnectionAuthType = "UserEntraToken"
	ConnectionAuthTypeUsernamePassword       ConnectionAuthType = "UsernamePassword"
)

func PossibleValuesForConnectionAuthType() []string {
	return []string{
		string(ConnectionAuthTypeAAD),
		string(ConnectionAuthTypeAccessKey),
		string(ConnectionAuthTypeAccountKey),
		string(ConnectionAuthTypeAccountManagedIdentity),
		string(ConnectionAuthTypeAgentUserImpersonation),
		string(ConnectionAuthTypeAgenticIdentityToken),
		string(ConnectionAuthTypeAgenticUser),
		string(ConnectionAuthTypeApiKey),
		string(ConnectionAuthTypeCustomKeys),
		string(ConnectionAuthTypeDelegatedSAS),
		string(ConnectionAuthTypeManagedIdentity),
		string(ConnectionAuthTypeNone),
		string(ConnectionAuthTypeOAuthTwo),
		string(ConnectionAuthTypePAT),
		string(ConnectionAuthTypeProjectManagedIdentity),
		string(ConnectionAuthTypeSAS),
		string(ConnectionAuthTypeServicePrincipal),
		string(ConnectionAuthTypeUserEntraToken),
		string(ConnectionAuthTypeUsernamePassword),
	}
}

func (s *ConnectionAuthType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionAuthType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionAuthType(input string) (*ConnectionAuthType, error) {
	vals := map[string]ConnectionAuthType{
		"aad":                    ConnectionAuthTypeAAD,
		"accesskey":              ConnectionAuthTypeAccessKey,
		"accountkey":             ConnectionAuthTypeAccountKey,
		"accountmanagedidentity": ConnectionAuthTypeAccountManagedIdentity,
		"agentuserimpersonation": ConnectionAuthTypeAgentUserImpersonation,
		"agenticidentitytoken":   ConnectionAuthTypeAgenticIdentityToken,
		"agenticuser":            ConnectionAuthTypeAgenticUser,
		"apikey":                 ConnectionAuthTypeApiKey,
		"customkeys":             ConnectionAuthTypeCustomKeys,
		"delegatedsas":           ConnectionAuthTypeDelegatedSAS,
		"managedidentity":        ConnectionAuthTypeManagedIdentity,
		"none":                   ConnectionAuthTypeNone,
		"oauth2":                 ConnectionAuthTypeOAuthTwo,
		"pat":                    ConnectionAuthTypePAT,
		"projectmanagedidentity": ConnectionAuthTypeProjectManagedIdentity,
		"sas":                    ConnectionAuthTypeSAS,
		"serviceprincipal":       ConnectionAuthTypeServicePrincipal,
		"userentratoken":         ConnectionAuthTypeUserEntraToken,
		"usernamepassword":       ConnectionAuthTypeUsernamePassword,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionAuthType(input)
	return &out, nil
}

type ConnectionCategory string

const (
	ConnectionCategoryADLSGenTwo                   ConnectionCategory = "ADLSGen2"
	ConnectionCategoryAIServices                   ConnectionCategory = "AIServices"
	ConnectionCategoryAmazonMws                    ConnectionCategory = "AmazonMws"
	ConnectionCategoryAmazonRdsForOracle           ConnectionCategory = "AmazonRdsForOracle"
	ConnectionCategoryAmazonRdsForSqlServer        ConnectionCategory = "AmazonRdsForSqlServer"
	ConnectionCategoryAmazonRedshift               ConnectionCategory = "AmazonRedshift"
	ConnectionCategoryAmazonSThreeCompatible       ConnectionCategory = "AmazonS3Compatible"
	ConnectionCategoryApiKey                       ConnectionCategory = "ApiKey"
	ConnectionCategoryApiManagement                ConnectionCategory = "ApiManagement"
	ConnectionCategoryAppConfig                    ConnectionCategory = "AppConfig"
	ConnectionCategoryAppInsights                  ConnectionCategory = "AppInsights"
	ConnectionCategoryAzureBlob                    ConnectionCategory = "AzureBlob"
	ConnectionCategoryAzureContainerAppEnvironment ConnectionCategory = "AzureContainerAppEnvironment"
	ConnectionCategoryAzureDataExplorer            ConnectionCategory = "AzureDataExplorer"
	ConnectionCategoryAzureDatabricksDeltaLake     ConnectionCategory = "AzureDatabricksDeltaLake"
	ConnectionCategoryAzureKeyVault                ConnectionCategory = "AzureKeyVault"
	ConnectionCategoryAzureMariaDb                 ConnectionCategory = "AzureMariaDb"
	ConnectionCategoryAzureMySqlDb                 ConnectionCategory = "AzureMySqlDb"
	ConnectionCategoryAzureOneLake                 ConnectionCategory = "AzureOneLake"
	ConnectionCategoryAzureOpenAI                  ConnectionCategory = "AzureOpenAI"
	ConnectionCategoryAzurePostgresDb              ConnectionCategory = "AzurePostgresDb"
	ConnectionCategoryAzureSqlDb                   ConnectionCategory = "AzureSqlDb"
	ConnectionCategoryAzureSqlMi                   ConnectionCategory = "AzureSqlMi"
	ConnectionCategoryAzureStorageAccount          ConnectionCategory = "AzureStorageAccount"
	ConnectionCategoryAzureSynapseAnalytics        ConnectionCategory = "AzureSynapseAnalytics"
	ConnectionCategoryAzureTableStorage            ConnectionCategory = "AzureTableStorage"
	ConnectionCategoryBingLLMSearch                ConnectionCategory = "BingLLMSearch"
	ConnectionCategoryCassandra                    ConnectionCategory = "Cassandra"
	ConnectionCategoryCognitiveSearch              ConnectionCategory = "CognitiveSearch"
	ConnectionCategoryCognitiveService             ConnectionCategory = "CognitiveService"
	ConnectionCategoryConcur                       ConnectionCategory = "Concur"
	ConnectionCategoryContainerRegistry            ConnectionCategory = "ContainerRegistry"
	ConnectionCategoryCosmosDb                     ConnectionCategory = "CosmosDb"
	ConnectionCategoryCosmosDbMongoDbApi           ConnectionCategory = "CosmosDbMongoDbApi"
	ConnectionCategoryCouchbase                    ConnectionCategory = "Couchbase"
	ConnectionCategoryCustomKeys                   ConnectionCategory = "CustomKeys"
	ConnectionCategoryDatabricks                   ConnectionCategory = "Databricks"
	ConnectionCategoryDbTwo                        ConnectionCategory = "Db2"
	ConnectionCategoryDrill                        ConnectionCategory = "Drill"
	ConnectionCategoryDynamics                     ConnectionCategory = "Dynamics"
	ConnectionCategoryDynamicsAx                   ConnectionCategory = "DynamicsAx"
	ConnectionCategoryDynamicsCrm                  ConnectionCategory = "DynamicsCrm"
	ConnectionCategoryElasticsearch                ConnectionCategory = "Elasticsearch"
	ConnectionCategoryEloqua                       ConnectionCategory = "Eloqua"
	ConnectionCategoryFileServer                   ConnectionCategory = "FileServer"
	ConnectionCategoryFtpServer                    ConnectionCategory = "FtpServer"
	ConnectionCategoryGenericContainerRegistry     ConnectionCategory = "GenericContainerRegistry"
	ConnectionCategoryGenericHTTP                  ConnectionCategory = "GenericHttp"
	ConnectionCategoryGenericRest                  ConnectionCategory = "GenericRest"
	ConnectionCategoryGit                          ConnectionCategory = "Git"
	ConnectionCategoryGoogleAdWords                ConnectionCategory = "GoogleAdWords"
	ConnectionCategoryGoogleBigQuery               ConnectionCategory = "GoogleBigQuery"
	ConnectionCategoryGoogleCloudStorage           ConnectionCategory = "GoogleCloudStorage"
	ConnectionCategoryGreenplum                    ConnectionCategory = "Greenplum"
	ConnectionCategoryGroundingWithBingSearch      ConnectionCategory = "GroundingWithBingSearch"
	ConnectionCategoryGroundingWithCustomSearch    ConnectionCategory = "GroundingWithCustomSearch"
	ConnectionCategoryHbase                        ConnectionCategory = "Hbase"
	ConnectionCategoryHdfs                         ConnectionCategory = "Hdfs"
	ConnectionCategoryHive                         ConnectionCategory = "Hive"
	ConnectionCategoryHubspot                      ConnectionCategory = "Hubspot"
	ConnectionCategoryImpala                       ConnectionCategory = "Impala"
	ConnectionCategoryInformix                     ConnectionCategory = "Informix"
	ConnectionCategoryJira                         ConnectionCategory = "Jira"
	ConnectionCategoryMagento                      ConnectionCategory = "Magento"
	ConnectionCategoryManagedOnlineEndpoint        ConnectionCategory = "ManagedOnlineEndpoint"
	ConnectionCategoryMariaDb                      ConnectionCategory = "MariaDb"
	ConnectionCategoryMarketo                      ConnectionCategory = "Marketo"
	ConnectionCategoryMicrosoftAccess              ConnectionCategory = "MicrosoftAccess"
	ConnectionCategoryMicrosoftFabric              ConnectionCategory = "MicrosoftFabric"
	ConnectionCategoryModelGateway                 ConnectionCategory = "ModelGateway"
	ConnectionCategoryMongoDbAtlas                 ConnectionCategory = "MongoDbAtlas"
	ConnectionCategoryMongoDbVTwo                  ConnectionCategory = "MongoDbV2"
	ConnectionCategoryMySql                        ConnectionCategory = "MySql"
	ConnectionCategoryNetezza                      ConnectionCategory = "Netezza"
	ConnectionCategoryODataRest                    ConnectionCategory = "ODataRest"
	ConnectionCategoryOdbc                         ConnectionCategory = "Odbc"
	ConnectionCategoryOfficeThreeSixFive           ConnectionCategory = "Office365"
	ConnectionCategoryOpenAI                       ConnectionCategory = "OpenAI"
	ConnectionCategoryOracle                       ConnectionCategory = "Oracle"
	ConnectionCategoryOracleCloudStorage           ConnectionCategory = "OracleCloudStorage"
	ConnectionCategoryOracleServiceCloud           ConnectionCategory = "OracleServiceCloud"
	ConnectionCategoryPayPal                       ConnectionCategory = "PayPal"
	ConnectionCategoryPhoenix                      ConnectionCategory = "Phoenix"
	ConnectionCategoryPinecone                     ConnectionCategory = "Pinecone"
	ConnectionCategoryPostgreSql                   ConnectionCategory = "PostgreSql"
	ConnectionCategoryPowerPlatformEnvironment     ConnectionCategory = "PowerPlatformEnvironment"
	ConnectionCategoryPresto                       ConnectionCategory = "Presto"
	ConnectionCategoryPythonFeed                   ConnectionCategory = "PythonFeed"
	ConnectionCategoryQuickBooks                   ConnectionCategory = "QuickBooks"
	ConnectionCategoryRedis                        ConnectionCategory = "Redis"
	ConnectionCategoryRemoteATwoA                  ConnectionCategory = "RemoteA2A"
	ConnectionCategoryRemoteTool                   ConnectionCategory = "RemoteTool"
	ConnectionCategoryResponsys                    ConnectionCategory = "Responsys"
	ConnectionCategorySThree                       ConnectionCategory = "S3"
	ConnectionCategorySalesforce                   ConnectionCategory = "Salesforce"
	ConnectionCategorySalesforceMarketingCloud     ConnectionCategory = "SalesforceMarketingCloud"
	ConnectionCategorySalesforceServiceCloud       ConnectionCategory = "SalesforceServiceCloud"
	ConnectionCategorySapBw                        ConnectionCategory = "SapBw"
	ConnectionCategorySapCloudForCustomer          ConnectionCategory = "SapCloudForCustomer"
	ConnectionCategorySapEcc                       ConnectionCategory = "SapEcc"
	ConnectionCategorySapHana                      ConnectionCategory = "SapHana"
	ConnectionCategorySapOpenHub                   ConnectionCategory = "SapOpenHub"
	ConnectionCategorySapTable                     ConnectionCategory = "SapTable"
	ConnectionCategorySerp                         ConnectionCategory = "Serp"
	ConnectionCategoryServerless                   ConnectionCategory = "Serverless"
	ConnectionCategoryServiceNow                   ConnectionCategory = "ServiceNow"
	ConnectionCategorySftp                         ConnectionCategory = "Sftp"
	ConnectionCategorySharePointOnlineList         ConnectionCategory = "SharePointOnlineList"
	ConnectionCategorySharepoint                   ConnectionCategory = "Sharepoint"
	ConnectionCategoryShopify                      ConnectionCategory = "Shopify"
	ConnectionCategorySnowflake                    ConnectionCategory = "Snowflake"
	ConnectionCategorySpark                        ConnectionCategory = "Spark"
	ConnectionCategorySqlServer                    ConnectionCategory = "SqlServer"
	ConnectionCategorySquare                       ConnectionCategory = "Square"
	ConnectionCategorySybase                       ConnectionCategory = "Sybase"
	ConnectionCategoryTeradata                     ConnectionCategory = "Teradata"
	ConnectionCategoryVertica                      ConnectionCategory = "Vertica"
	ConnectionCategoryWebTable                     ConnectionCategory = "WebTable"
	ConnectionCategoryXero                         ConnectionCategory = "Xero"
	ConnectionCategoryZoho                         ConnectionCategory = "Zoho"
)

func PossibleValuesForConnectionCategory() []string {
	return []string{
		string(ConnectionCategoryADLSGenTwo),
		string(ConnectionCategoryAIServices),
		string(ConnectionCategoryAmazonMws),
		string(ConnectionCategoryAmazonRdsForOracle),
		string(ConnectionCategoryAmazonRdsForSqlServer),
		string(ConnectionCategoryAmazonRedshift),
		string(ConnectionCategoryAmazonSThreeCompatible),
		string(ConnectionCategoryApiKey),
		string(ConnectionCategoryApiManagement),
		string(ConnectionCategoryAppConfig),
		string(ConnectionCategoryAppInsights),
		string(ConnectionCategoryAzureBlob),
		string(ConnectionCategoryAzureContainerAppEnvironment),
		string(ConnectionCategoryAzureDataExplorer),
		string(ConnectionCategoryAzureDatabricksDeltaLake),
		string(ConnectionCategoryAzureKeyVault),
		string(ConnectionCategoryAzureMariaDb),
		string(ConnectionCategoryAzureMySqlDb),
		string(ConnectionCategoryAzureOneLake),
		string(ConnectionCategoryAzureOpenAI),
		string(ConnectionCategoryAzurePostgresDb),
		string(ConnectionCategoryAzureSqlDb),
		string(ConnectionCategoryAzureSqlMi),
		string(ConnectionCategoryAzureStorageAccount),
		string(ConnectionCategoryAzureSynapseAnalytics),
		string(ConnectionCategoryAzureTableStorage),
		string(ConnectionCategoryBingLLMSearch),
		string(ConnectionCategoryCassandra),
		string(ConnectionCategoryCognitiveSearch),
		string(ConnectionCategoryCognitiveService),
		string(ConnectionCategoryConcur),
		string(ConnectionCategoryContainerRegistry),
		string(ConnectionCategoryCosmosDb),
		string(ConnectionCategoryCosmosDbMongoDbApi),
		string(ConnectionCategoryCouchbase),
		string(ConnectionCategoryCustomKeys),
		string(ConnectionCategoryDatabricks),
		string(ConnectionCategoryDbTwo),
		string(ConnectionCategoryDrill),
		string(ConnectionCategoryDynamics),
		string(ConnectionCategoryDynamicsAx),
		string(ConnectionCategoryDynamicsCrm),
		string(ConnectionCategoryElasticsearch),
		string(ConnectionCategoryEloqua),
		string(ConnectionCategoryFileServer),
		string(ConnectionCategoryFtpServer),
		string(ConnectionCategoryGenericContainerRegistry),
		string(ConnectionCategoryGenericHTTP),
		string(ConnectionCategoryGenericRest),
		string(ConnectionCategoryGit),
		string(ConnectionCategoryGoogleAdWords),
		string(ConnectionCategoryGoogleBigQuery),
		string(ConnectionCategoryGoogleCloudStorage),
		string(ConnectionCategoryGreenplum),
		string(ConnectionCategoryGroundingWithBingSearch),
		string(ConnectionCategoryGroundingWithCustomSearch),
		string(ConnectionCategoryHbase),
		string(ConnectionCategoryHdfs),
		string(ConnectionCategoryHive),
		string(ConnectionCategoryHubspot),
		string(ConnectionCategoryImpala),
		string(ConnectionCategoryInformix),
		string(ConnectionCategoryJira),
		string(ConnectionCategoryMagento),
		string(ConnectionCategoryManagedOnlineEndpoint),
		string(ConnectionCategoryMariaDb),
		string(ConnectionCategoryMarketo),
		string(ConnectionCategoryMicrosoftAccess),
		string(ConnectionCategoryMicrosoftFabric),
		string(ConnectionCategoryModelGateway),
		string(ConnectionCategoryMongoDbAtlas),
		string(ConnectionCategoryMongoDbVTwo),
		string(ConnectionCategoryMySql),
		string(ConnectionCategoryNetezza),
		string(ConnectionCategoryODataRest),
		string(ConnectionCategoryOdbc),
		string(ConnectionCategoryOfficeThreeSixFive),
		string(ConnectionCategoryOpenAI),
		string(ConnectionCategoryOracle),
		string(ConnectionCategoryOracleCloudStorage),
		string(ConnectionCategoryOracleServiceCloud),
		string(ConnectionCategoryPayPal),
		string(ConnectionCategoryPhoenix),
		string(ConnectionCategoryPinecone),
		string(ConnectionCategoryPostgreSql),
		string(ConnectionCategoryPowerPlatformEnvironment),
		string(ConnectionCategoryPresto),
		string(ConnectionCategoryPythonFeed),
		string(ConnectionCategoryQuickBooks),
		string(ConnectionCategoryRedis),
		string(ConnectionCategoryRemoteATwoA),
		string(ConnectionCategoryRemoteTool),
		string(ConnectionCategoryResponsys),
		string(ConnectionCategorySThree),
		string(ConnectionCategorySalesforce),
		string(ConnectionCategorySalesforceMarketingCloud),
		string(ConnectionCategorySalesforceServiceCloud),
		string(ConnectionCategorySapBw),
		string(ConnectionCategorySapCloudForCustomer),
		string(ConnectionCategorySapEcc),
		string(ConnectionCategorySapHana),
		string(ConnectionCategorySapOpenHub),
		string(ConnectionCategorySapTable),
		string(ConnectionCategorySerp),
		string(ConnectionCategoryServerless),
		string(ConnectionCategoryServiceNow),
		string(ConnectionCategorySftp),
		string(ConnectionCategorySharePointOnlineList),
		string(ConnectionCategorySharepoint),
		string(ConnectionCategoryShopify),
		string(ConnectionCategorySnowflake),
		string(ConnectionCategorySpark),
		string(ConnectionCategorySqlServer),
		string(ConnectionCategorySquare),
		string(ConnectionCategorySybase),
		string(ConnectionCategoryTeradata),
		string(ConnectionCategoryVertica),
		string(ConnectionCategoryWebTable),
		string(ConnectionCategoryXero),
		string(ConnectionCategoryZoho),
	}
}

func (s *ConnectionCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionCategory(input string) (*ConnectionCategory, error) {
	vals := map[string]ConnectionCategory{
		"adlsgen2":                     ConnectionCategoryADLSGenTwo,
		"aiservices":                   ConnectionCategoryAIServices,
		"amazonmws":                    ConnectionCategoryAmazonMws,
		"amazonrdsfororacle":           ConnectionCategoryAmazonRdsForOracle,
		"amazonrdsforsqlserver":        ConnectionCategoryAmazonRdsForSqlServer,
		"amazonredshift":               ConnectionCategoryAmazonRedshift,
		"amazons3compatible":           ConnectionCategoryAmazonSThreeCompatible,
		"apikey":                       ConnectionCategoryApiKey,
		"apimanagement":                ConnectionCategoryApiManagement,
		"appconfig":                    ConnectionCategoryAppConfig,
		"appinsights":                  ConnectionCategoryAppInsights,
		"azureblob":                    ConnectionCategoryAzureBlob,
		"azurecontainerappenvironment": ConnectionCategoryAzureContainerAppEnvironment,
		"azuredataexplorer":            ConnectionCategoryAzureDataExplorer,
		"azuredatabricksdeltalake":     ConnectionCategoryAzureDatabricksDeltaLake,
		"azurekeyvault":                ConnectionCategoryAzureKeyVault,
		"azuremariadb":                 ConnectionCategoryAzureMariaDb,
		"azuremysqldb":                 ConnectionCategoryAzureMySqlDb,
		"azureonelake":                 ConnectionCategoryAzureOneLake,
		"azureopenai":                  ConnectionCategoryAzureOpenAI,
		"azurepostgresdb":              ConnectionCategoryAzurePostgresDb,
		"azuresqldb":                   ConnectionCategoryAzureSqlDb,
		"azuresqlmi":                   ConnectionCategoryAzureSqlMi,
		"azurestorageaccount":          ConnectionCategoryAzureStorageAccount,
		"azuresynapseanalytics":        ConnectionCategoryAzureSynapseAnalytics,
		"azuretablestorage":            ConnectionCategoryAzureTableStorage,
		"bingllmsearch":                ConnectionCategoryBingLLMSearch,
		"cassandra":                    ConnectionCategoryCassandra,
		"cognitivesearch":              ConnectionCategoryCognitiveSearch,
		"cognitiveservice":             ConnectionCategoryCognitiveService,
		"concur":                       ConnectionCategoryConcur,
		"containerregistry":            ConnectionCategoryContainerRegistry,
		"cosmosdb":                     ConnectionCategoryCosmosDb,
		"cosmosdbmongodbapi":           ConnectionCategoryCosmosDbMongoDbApi,
		"couchbase":                    ConnectionCategoryCouchbase,
		"customkeys":                   ConnectionCategoryCustomKeys,
		"databricks":                   ConnectionCategoryDatabricks,
		"db2":                          ConnectionCategoryDbTwo,
		"drill":                        ConnectionCategoryDrill,
		"dynamics":                     ConnectionCategoryDynamics,
		"dynamicsax":                   ConnectionCategoryDynamicsAx,
		"dynamicscrm":                  ConnectionCategoryDynamicsCrm,
		"elasticsearch":                ConnectionCategoryElasticsearch,
		"eloqua":                       ConnectionCategoryEloqua,
		"fileserver":                   ConnectionCategoryFileServer,
		"ftpserver":                    ConnectionCategoryFtpServer,
		"genericcontainerregistry":     ConnectionCategoryGenericContainerRegistry,
		"generichttp":                  ConnectionCategoryGenericHTTP,
		"genericrest":                  ConnectionCategoryGenericRest,
		"git":                          ConnectionCategoryGit,
		"googleadwords":                ConnectionCategoryGoogleAdWords,
		"googlebigquery":               ConnectionCategoryGoogleBigQuery,
		"googlecloudstorage":           ConnectionCategoryGoogleCloudStorage,
		"greenplum":                    ConnectionCategoryGreenplum,
		"groundingwithbingsearch":      ConnectionCategoryGroundingWithBingSearch,
		"groundingwithcustomsearch":    ConnectionCategoryGroundingWithCustomSearch,
		"hbase":                        ConnectionCategoryHbase,
		"hdfs":                         ConnectionCategoryHdfs,
		"hive":                         ConnectionCategoryHive,
		"hubspot":                      ConnectionCategoryHubspot,
		"impala":                       ConnectionCategoryImpala,
		"informix":                     ConnectionCategoryInformix,
		"jira":                         ConnectionCategoryJira,
		"magento":                      ConnectionCategoryMagento,
		"managedonlineendpoint":        ConnectionCategoryManagedOnlineEndpoint,
		"mariadb":                      ConnectionCategoryMariaDb,
		"marketo":                      ConnectionCategoryMarketo,
		"microsoftaccess":              ConnectionCategoryMicrosoftAccess,
		"microsoftfabric":              ConnectionCategoryMicrosoftFabric,
		"modelgateway":                 ConnectionCategoryModelGateway,
		"mongodbatlas":                 ConnectionCategoryMongoDbAtlas,
		"mongodbv2":                    ConnectionCategoryMongoDbVTwo,
		"mysql":                        ConnectionCategoryMySql,
		"netezza":                      ConnectionCategoryNetezza,
		"odatarest":                    ConnectionCategoryODataRest,
		"odbc":                         ConnectionCategoryOdbc,
		"office365":                    ConnectionCategoryOfficeThreeSixFive,
		"openai":                       ConnectionCategoryOpenAI,
		"oracle":                       ConnectionCategoryOracle,
		"oraclecloudstorage":           ConnectionCategoryOracleCloudStorage,
		"oracleservicecloud":           ConnectionCategoryOracleServiceCloud,
		"paypal":                       ConnectionCategoryPayPal,
		"phoenix":                      ConnectionCategoryPhoenix,
		"pinecone":                     ConnectionCategoryPinecone,
		"postgresql":                   ConnectionCategoryPostgreSql,
		"powerplatformenvironment":     ConnectionCategoryPowerPlatformEnvironment,
		"presto":                       ConnectionCategoryPresto,
		"pythonfeed":                   ConnectionCategoryPythonFeed,
		"quickbooks":                   ConnectionCategoryQuickBooks,
		"redis":                        ConnectionCategoryRedis,
		"remotea2a":                    ConnectionCategoryRemoteATwoA,
		"remotetool":                   ConnectionCategoryRemoteTool,
		"responsys":                    ConnectionCategoryResponsys,
		"s3":                           ConnectionCategorySThree,
		"salesforce":                   ConnectionCategorySalesforce,
		"salesforcemarketingcloud":     ConnectionCategorySalesforceMarketingCloud,
		"salesforceservicecloud":       ConnectionCategorySalesforceServiceCloud,
		"sapbw":                        ConnectionCategorySapBw,
		"sapcloudforcustomer":          ConnectionCategorySapCloudForCustomer,
		"sapecc":                       ConnectionCategorySapEcc,
		"saphana":                      ConnectionCategorySapHana,
		"sapopenhub":                   ConnectionCategorySapOpenHub,
		"saptable":                     ConnectionCategorySapTable,
		"serp":                         ConnectionCategorySerp,
		"serverless":                   ConnectionCategoryServerless,
		"servicenow":                   ConnectionCategoryServiceNow,
		"sftp":                         ConnectionCategorySftp,
		"sharepointonlinelist":         ConnectionCategorySharePointOnlineList,
		"sharepoint":                   ConnectionCategorySharepoint,
		"shopify":                      ConnectionCategoryShopify,
		"snowflake":                    ConnectionCategorySnowflake,
		"spark":                        ConnectionCategorySpark,
		"sqlserver":                    ConnectionCategorySqlServer,
		"square":                       ConnectionCategorySquare,
		"sybase":                       ConnectionCategorySybase,
		"teradata":                     ConnectionCategoryTeradata,
		"vertica":                      ConnectionCategoryVertica,
		"webtable":                     ConnectionCategoryWebTable,
		"xero":                         ConnectionCategoryXero,
		"zoho":                         ConnectionCategoryZoho,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionCategory(input)
	return &out, nil
}

type ConnectionGroup string

const (
	ConnectionGroupAzure           ConnectionGroup = "Azure"
	ConnectionGroupAzureAI         ConnectionGroup = "AzureAI"
	ConnectionGroupDatabase        ConnectionGroup = "Database"
	ConnectionGroupFile            ConnectionGroup = "File"
	ConnectionGroupGenericProtocol ConnectionGroup = "GenericProtocol"
	ConnectionGroupNoSQL           ConnectionGroup = "NoSQL"
	ConnectionGroupServicesAndApps ConnectionGroup = "ServicesAndApps"
)

func PossibleValuesForConnectionGroup() []string {
	return []string{
		string(ConnectionGroupAzure),
		string(ConnectionGroupAzureAI),
		string(ConnectionGroupDatabase),
		string(ConnectionGroupFile),
		string(ConnectionGroupGenericProtocol),
		string(ConnectionGroupNoSQL),
		string(ConnectionGroupServicesAndApps),
	}
}

func (s *ConnectionGroup) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionGroup(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionGroup(input string) (*ConnectionGroup, error) {
	vals := map[string]ConnectionGroup{
		"azure":           ConnectionGroupAzure,
		"azureai":         ConnectionGroupAzureAI,
		"database":        ConnectionGroupDatabase,
		"file":            ConnectionGroupFile,
		"genericprotocol": ConnectionGroupGenericProtocol,
		"nosql":           ConnectionGroupNoSQL,
		"servicesandapps": ConnectionGroupServicesAndApps,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionGroup(input)
	return &out, nil
}

type ManagedPERequirement string

const (
	ManagedPERequirementNotApplicable ManagedPERequirement = "NotApplicable"
	ManagedPERequirementNotRequired   ManagedPERequirement = "NotRequired"
	ManagedPERequirementRequired      ManagedPERequirement = "Required"
)

func PossibleValuesForManagedPERequirement() []string {
	return []string{
		string(ManagedPERequirementNotApplicable),
		string(ManagedPERequirementNotRequired),
		string(ManagedPERequirementRequired),
	}
}

func (s *ManagedPERequirement) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedPERequirement(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedPERequirement(input string) (*ManagedPERequirement, error) {
	vals := map[string]ManagedPERequirement{
		"notapplicable": ManagedPERequirementNotApplicable,
		"notrequired":   ManagedPERequirementNotRequired,
		"required":      ManagedPERequirementRequired,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedPERequirement(input)
	return &out, nil
}

type ManagedPEStatus string

const (
	ManagedPEStatusActive        ManagedPEStatus = "Active"
	ManagedPEStatusInactive      ManagedPEStatus = "Inactive"
	ManagedPEStatusNotApplicable ManagedPEStatus = "NotApplicable"
)

func PossibleValuesForManagedPEStatus() []string {
	return []string{
		string(ManagedPEStatusActive),
		string(ManagedPEStatusInactive),
		string(ManagedPEStatusNotApplicable),
	}
}

func (s *ManagedPEStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedPEStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedPEStatus(input string) (*ManagedPEStatus, error) {
	vals := map[string]ManagedPEStatus{
		"active":        ManagedPEStatusActive,
		"inactive":      ManagedPEStatusInactive,
		"notapplicable": ManagedPEStatusNotApplicable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedPEStatus(input)
	return &out, nil
}
