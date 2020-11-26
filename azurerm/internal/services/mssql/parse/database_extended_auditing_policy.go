package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabaseExtendedAuditingPolicyId struct {
	DatabaseName  string
	ServerName    string
	ResourceGroup string
}

func DatabaseExtendedAuditingPolicyID(input string) (*DatabaseExtendedAuditingPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql Database Extended Auditing Policy %q: %+v", input, err)
	}

	sqlDatabaseExtendedAuditingPolicyId := DatabaseExtendedAuditingPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if sqlDatabaseExtendedAuditingPolicyId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if sqlDatabaseExtendedAuditingPolicyId.DatabaseName, err = id.PopSegment("databases"); err != nil {
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
