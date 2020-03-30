package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StorageImportExportJobId struct {
	ResourceGroup string
	Name          string
}

func StorageImportExportJobID(input string) (*StorageImportExportJobId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse StorageImportExportJob ID %q: %+v", input, err)
	}

	storageImportExportJob := StorageImportExportJobId{
		ResourceGroup: id.ResourceGroup,
	}

	if storageImportExportJob.Name, err = id.PopSegment("jobs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &storageImportExportJob, nil
}
