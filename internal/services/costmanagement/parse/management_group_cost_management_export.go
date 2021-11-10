package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type ManagementGroupCostManagementExportId struct {
	ManagementGroupName string
	ExportName          string
}

func NewManagementGroupCostManagementExportID(managementGroupName, exportName string) ManagementGroupCostManagementExportId {
	return ManagementGroupCostManagementExportId{
		ManagementGroupName: managementGroupName,
		ExportName:          exportName,
	}
}

func (id ManagementGroupCostManagementExportId) String() string {
	segments := []string{
		fmt.Sprintf("Export Name %q", id.ExportName),
		fmt.Sprintf("Management Group Name %q", id.ManagementGroupName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Management Group Cost Management Export", segmentsStr)
}

func (id ManagementGroupCostManagementExportId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.CostManagement/exports/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupName, id.ExportName)
}

// ManagementGroupCostManagementExportID parses a ManagementGroupCostManagementExport ID into an ManagementGroupCostManagementExportId struct
func ManagementGroupCostManagementExportID(input string) (*ManagementGroupCostManagementExportId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagementGroupCostManagementExportId{}

	if resourceId.ManagementGroupName, err = id.PopSegment("managementGroups"); err != nil {
		return nil, err
	}
	if resourceId.ExportName, err = id.PopSegment("exports"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
