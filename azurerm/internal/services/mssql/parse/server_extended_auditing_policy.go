package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MsSqlServerExtendedAuditingPolicyId struct {
	MsSqlServer   string
	ResourceGroup string
}

func MssqlServerExtendedAuditingPolicyID(input string) (*MsSqlServerExtendedAuditingPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql Server Extended Auditing Policy %q: %+v", input, err)
	}

	sqlServerExtendedAuditingPolicyId := MsSqlServerExtendedAuditingPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if sqlServerExtendedAuditingPolicyId.MsSqlServer, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if _, err = id.PopSegment("extendedAuditingSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &sqlServerExtendedAuditingPolicyId, nil
}
