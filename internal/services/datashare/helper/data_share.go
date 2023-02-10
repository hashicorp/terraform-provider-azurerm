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
	if d, ok := ds.(dataset.ADLSGen2FileDataSet); ok {
		return d.Id
	}
	if d, ok := ds.(dataset.ADLSGen2FolderDataSet); ok {
		return d.Id
	}
	if d, ok := ds.(dataset.ADLSGen2FileSystemDataSet); ok {
		return d.Id
	}
	if d, ok := ds.(dataset.KustoClusterDataSet); ok {
		return d.Id
	}
	if d, ok := ds.(dataset.KustoDatabaseDataSet); ok {
		return d.Id
	}

	return nil
}
