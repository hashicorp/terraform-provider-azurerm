package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MsSqlDatabaseExtendedAuditingPolicyId struct {
	MsDBName      string
	MsSqlServer   string
	ResourceGroup string
}

func MssqlDatabaseExtendedAuditingPolicyID(input string) (*MsSqlDatabaseExtendedAuditingPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql Database Extended Auditing Policy %q: %+v", input, err)
	}

	sqlDatabaseExtendedAuditingPolicyId := MsSqlDatabaseExtendedAuditingPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if sqlDatabaseExtendedAuditingPolicyId.MsSqlServer, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if sqlDatabaseExtendedAuditingPolicyId.MsDBName, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}

	if _, err = id.PopSegment("extendedAuditingSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &sqlDatabaseExtendedAuditingPolicyId, nil
}
