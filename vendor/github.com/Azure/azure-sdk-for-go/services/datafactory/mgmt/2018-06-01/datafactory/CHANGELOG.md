# Change History

## Breaking Changes

### Removed Constants

1. TypeBasicDatasetCompression.TypeBasicDatasetCompressionTypeBZip2
1. TypeBasicDatasetCompression.TypeBasicDatasetCompressionTypeDatasetCompression
1. TypeBasicDatasetCompression.TypeBasicDatasetCompressionTypeDeflate
1. TypeBasicDatasetCompression.TypeBasicDatasetCompressionTypeGZip
1. TypeBasicDatasetCompression.TypeBasicDatasetCompressionTypeTar
1. TypeBasicDatasetCompression.TypeBasicDatasetCompressionTypeTarGZip
1. TypeBasicDatasetCompression.TypeBasicDatasetCompressionTypeZipDeflate

### Removed Funcs

1. *DatasetBZip2Compression.UnmarshalJSON([]byte) error
1. *DatasetDeflateCompression.UnmarshalJSON([]byte) error
1. *DatasetGZipCompression.UnmarshalJSON([]byte) error
1. *DatasetTarCompression.UnmarshalJSON([]byte) error
1. *DatasetTarGZipCompression.UnmarshalJSON([]byte) error
1. *DatasetZipDeflateCompression.UnmarshalJSON([]byte) error
1. DatasetBZip2Compression.AsBasicDatasetCompression() (BasicDatasetCompression, bool)
1. DatasetBZip2Compression.AsDatasetBZip2Compression() (*DatasetBZip2Compression, bool)
1. DatasetBZip2Compression.AsDatasetCompression() (*DatasetCompression, bool)
1. DatasetBZip2Compression.AsDatasetDeflateCompression() (*DatasetDeflateCompression, bool)
1. DatasetBZip2Compression.AsDatasetGZipCompression() (*DatasetGZipCompression, bool)
1. DatasetBZip2Compression.AsDatasetTarCompression() (*DatasetTarCompression, bool)
1. DatasetBZip2Compression.AsDatasetTarGZipCompression() (*DatasetTarGZipCompression, bool)
1. DatasetBZip2Compression.AsDatasetZipDeflateCompression() (*DatasetZipDeflateCompression, bool)
1. DatasetBZip2Compression.MarshalJSON() ([]byte, error)
1. DatasetCompression.AsBasicDatasetCompression() (BasicDatasetCompression, bool)
1. DatasetCompression.AsDatasetBZip2Compression() (*DatasetBZip2Compression, bool)
1. DatasetCompression.AsDatasetCompression() (*DatasetCompression, bool)
1. DatasetCompression.AsDatasetDeflateCompression() (*DatasetDeflateCompression, bool)
1. DatasetCompression.AsDatasetGZipCompression() (*DatasetGZipCompression, bool)
1. DatasetCompression.AsDatasetTarCompression() (*DatasetTarCompression, bool)
1. DatasetCompression.AsDatasetTarGZipCompression() (*DatasetTarGZipCompression, bool)
1. DatasetCompression.AsDatasetZipDeflateCompression() (*DatasetZipDeflateCompression, bool)
1. DatasetDeflateCompression.AsBasicDatasetCompression() (BasicDatasetCompression, bool)
1. DatasetDeflateCompression.AsDatasetBZip2Compression() (*DatasetBZip2Compression, bool)
1. DatasetDeflateCompression.AsDatasetCompression() (*DatasetCompression, bool)
1. DatasetDeflateCompression.AsDatasetDeflateCompression() (*DatasetDeflateCompression, bool)
1. DatasetDeflateCompression.AsDatasetGZipCompression() (*DatasetGZipCompression, bool)
1. DatasetDeflateCompression.AsDatasetTarCompression() (*DatasetTarCompression, bool)
1. DatasetDeflateCompression.AsDatasetTarGZipCompression() (*DatasetTarGZipCompression, bool)
1. DatasetDeflateCompression.AsDatasetZipDeflateCompression() (*DatasetZipDeflateCompression, bool)
1. DatasetDeflateCompression.MarshalJSON() ([]byte, error)
1. DatasetGZipCompression.AsBasicDatasetCompression() (BasicDatasetCompression, bool)
1. DatasetGZipCompression.AsDatasetBZip2Compression() (*DatasetBZip2Compression, bool)
1. DatasetGZipCompression.AsDatasetCompression() (*DatasetCompression, bool)
1. DatasetGZipCompression.AsDatasetDeflateCompression() (*DatasetDeflateCompression, bool)
1. DatasetGZipCompression.AsDatasetGZipCompression() (*DatasetGZipCompression, bool)
1. DatasetGZipCompression.AsDatasetTarCompression() (*DatasetTarCompression, bool)
1. DatasetGZipCompression.AsDatasetTarGZipCompression() (*DatasetTarGZipCompression, bool)
1. DatasetGZipCompression.AsDatasetZipDeflateCompression() (*DatasetZipDeflateCompression, bool)
1. DatasetGZipCompression.MarshalJSON() ([]byte, error)
1. DatasetTarCompression.AsBasicDatasetCompression() (BasicDatasetCompression, bool)
1. DatasetTarCompression.AsDatasetBZip2Compression() (*DatasetBZip2Compression, bool)
1. DatasetTarCompression.AsDatasetCompression() (*DatasetCompression, bool)
1. DatasetTarCompression.AsDatasetDeflateCompression() (*DatasetDeflateCompression, bool)
1. DatasetTarCompression.AsDatasetGZipCompression() (*DatasetGZipCompression, bool)
1. DatasetTarCompression.AsDatasetTarCompression() (*DatasetTarCompression, bool)
1. DatasetTarCompression.AsDatasetTarGZipCompression() (*DatasetTarGZipCompression, bool)
1. DatasetTarCompression.AsDatasetZipDeflateCompression() (*DatasetZipDeflateCompression, bool)
1. DatasetTarCompression.MarshalJSON() ([]byte, error)
1. DatasetTarGZipCompression.AsBasicDatasetCompression() (BasicDatasetCompression, bool)
1. DatasetTarGZipCompression.AsDatasetBZip2Compression() (*DatasetBZip2Compression, bool)
1. DatasetTarGZipCompression.AsDatasetCompression() (*DatasetCompression, bool)
1. DatasetTarGZipCompression.AsDatasetDeflateCompression() (*DatasetDeflateCompression, bool)
1. DatasetTarGZipCompression.AsDatasetGZipCompression() (*DatasetGZipCompression, bool)
1. DatasetTarGZipCompression.AsDatasetTarCompression() (*DatasetTarCompression, bool)
1. DatasetTarGZipCompression.AsDatasetTarGZipCompression() (*DatasetTarGZipCompression, bool)
1. DatasetTarGZipCompression.AsDatasetZipDeflateCompression() (*DatasetZipDeflateCompression, bool)
1. DatasetTarGZipCompression.MarshalJSON() ([]byte, error)
1. DatasetZipDeflateCompression.AsBasicDatasetCompression() (BasicDatasetCompression, bool)
1. DatasetZipDeflateCompression.AsDatasetBZip2Compression() (*DatasetBZip2Compression, bool)
1. DatasetZipDeflateCompression.AsDatasetCompression() (*DatasetCompression, bool)
1. DatasetZipDeflateCompression.AsDatasetDeflateCompression() (*DatasetDeflateCompression, bool)
1. DatasetZipDeflateCompression.AsDatasetGZipCompression() (*DatasetGZipCompression, bool)
1. DatasetZipDeflateCompression.AsDatasetTarCompression() (*DatasetTarCompression, bool)
1. DatasetZipDeflateCompression.AsDatasetTarGZipCompression() (*DatasetTarGZipCompression, bool)
1. DatasetZipDeflateCompression.AsDatasetZipDeflateCompression() (*DatasetZipDeflateCompression, bool)
1. DatasetZipDeflateCompression.MarshalJSON() ([]byte, error)
1. PossibleTypeBasicDatasetCompressionValues() []TypeBasicDatasetCompression

### Struct Changes

#### Removed Structs

1. DatasetBZip2Compression
1. DatasetDeflateCompression
1. DatasetGZipCompression
1. DatasetTarCompression
1. DatasetTarGZipCompression
1. DatasetZipDeflateCompression

### Signature Changes

#### Struct Fields

1. AmazonS3DatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. AzureBlobDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. AzureBlobFSDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. AzureDataLakeStoreDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. BinaryDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. DatasetCompression.Type changed type from TypeBasicDatasetCompression to interface{}
1. ExcelDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. FileShareDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. HTTPDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. JSONDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression
1. XMLDatasetTypeProperties.Compression changed type from BasicDatasetCompression to *DatasetCompression

## Additive Changes

### New Constants

1. TypeBasicDataFlow.TypeBasicDataFlowTypeFlowlet

### New Funcs

1. *Flowlet.UnmarshalJSON([]byte) error
1. DataFlow.AsFlowlet() (*Flowlet, bool)
1. DataFlowReference.MarshalJSON() ([]byte, error)
1. Flowlet.AsBasicDataFlow() (BasicDataFlow, bool)
1. Flowlet.AsDataFlow() (*DataFlow, bool)
1. Flowlet.AsFlowlet() (*Flowlet, bool)
1. Flowlet.AsMappingDataFlow() (*MappingDataFlow, bool)
1. Flowlet.AsWranglingDataFlow() (*WranglingDataFlow, bool)
1. Flowlet.MarshalJSON() ([]byte, error)
1. MappingDataFlow.AsFlowlet() (*Flowlet, bool)
1. WranglingDataFlow.AsFlowlet() (*Flowlet, bool)

### Struct Changes

#### New Structs

1. Flowlet
1. FlowletTypeProperties
1. PowerQuerySinkMapping

#### New Struct Fields

1. DataFlowDebugPackage.DataFlows
1. DataFlowReference.Parameters
1. DataFlowSink.Flowlet
1. DataFlowSource.Flowlet
1. DatasetCompression.Level
1. ExecutePowerQueryActivityTypeProperties.Queries
1. FactoryUpdateParameters.PublicNetworkAccess
1. FtpReadSettings.DisableChunking
1. MappingDataFlowTypeProperties.ScriptLines
1. PowerQuerySink.Flowlet
1. PowerQuerySource.Flowlet
1. PowerQueryTypeProperties.DocumentLocale
1. SftpReadSettings.DisableChunking
1. Transformation.Dataset
1. Transformation.Flowlet
1. Transformation.LinkedService
