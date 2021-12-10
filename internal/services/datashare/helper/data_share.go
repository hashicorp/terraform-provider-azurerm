package helper

import (
	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
)

func GetAzurermDataShareDataSetId(dataset datashare.BasicDataSet) *string {
	if dataset == nil {
		return nil
	}
	switch t := dataset.(type) {
	case datashare.BlobDataSet:
		return t.ID
	case datashare.BlobFolderDataSet:
		return t.ID
	case datashare.BlobContainerDataSet:
		return t.ID
	case datashare.ADLSGen2FileDataSet:
		return t.ID
	case datashare.ADLSGen2FolderDataSet:
		return t.ID
	case datashare.ADLSGen2FileSystemDataSet:
		return t.ID
	case datashare.ADLSGen1FolderDataSet:
		return t.ID
	case datashare.ADLSGen1FileDataSet:
		return t.ID
	case datashare.KustoClusterDataSet:
		return t.ID
	case datashare.KustoDatabaseDataSet:
		return t.ID
	case datashare.SQLDWTableDataSet:
		return t.ID
	case datashare.SQLDBTableDataSet:
		return t.ID
	default:
		return nil
	}
}
