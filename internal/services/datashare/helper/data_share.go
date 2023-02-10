package helper

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
)

func GetAzurermDataShareDataSetId(ds dataset.DataSet) *string {
	if ds == nil {
		return nil
	}

	if d, ok := ds.(dataset.BlobDataSet); ok {
		return d.Id
	}
	if d, ok := ds.(dataset.BlobFolderDataSet); ok {
		return d.Id
	}
	if d, ok := ds.(dataset.BlobContainerDataSet); ok {
		return d.Id
	}
	//switch t := dataset.(type) {
	//case dataset.:
	//	return t.ID
	//case datashare.BlobFolderDataSet:
	//	return t.ID
	//case datashare.BlobContainerDataSet:
	//	return t.ID
	//case datashare.ADLSGen2FileDataSet:
	//	return t.ID
	//case datashare.ADLSGen2FolderDataSet:
	//	return t.ID
	//case datashare.ADLSGen2FileSystemDataSet:
	//	return t.ID
	//case datashare.ADLSGen1FolderDataSet:
	//	return t.ID
	//case datashare.ADLSGen1FileDataSet:
	//	return t.ID
	//case datashare.KustoClusterDataSet:
	//	return t.ID
	//case datashare.KustoDatabaseDataSet:
	//	return t.ID
	//case datashare.SQLDWTableDataSet:
	//	return t.ID
	//case datashare.SQLDBTableDataSet:
	//	return t.ID
	//default:
	//	return nil
	//}
	return nil
}
