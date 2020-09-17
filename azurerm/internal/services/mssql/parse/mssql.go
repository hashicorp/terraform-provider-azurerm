package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MsSqlDatabaseId struct {
	Name          string
	MsSqlServer   string
	ResourceGroup string
}

type MsSqlServerId struct {
	Name          string
	ResourceGroup string
}

type MsSqlElasticPoolId struct {
	Name          string
	MsSqlServer   string
	ResourceGroup string
}

type MsSqlDatabaseExtendedAuditingPolicyId struct {
	MsDBName      string
	MsSqlServer   string
	ResourceGroup string
}

type MsSqlServerExtendedAuditingPolicyId struct {
	MsSqlServer   string
	ResourceGroup string
}

func NewMsSqlDatabaseID(resourceGroup, msSqlServer, name string) MsSqlDatabaseId {
	return MsSqlDatabaseId{
		ResourceGroup: resourceGroup,
		MsSqlServer:   msSqlServer,
		Name:          name,
	}
}

func (id MsSqlDatabaseId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/databases/%s", subscriptionId, id.ResourceGroup, id.MsSqlServer, id.Name)
}

func MsSqlDatabaseID(input string) (*MsSqlDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse MsSql Database ID %q: %+v", input, err)
	}

	database := MsSqlDatabaseId{
		ResourceGroup: id.ResourceGroup,
	}

	if database.MsSqlServer, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if database.Name, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &database, nil
}

func MsSqlServerID(input string) (*MsSqlServerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse MsSql Server ID %q: %+v", input, err)
	}

	server := MsSqlServerId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}

func MSSqlElasticPoolID(input string) (*MsSqlElasticPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse MsSql Elastic Pool ID %q: %+v", input, err)
	}

	elasticPool := MsSqlElasticPoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if elasticPool.MsSqlServer, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if elasticPool.Name, err = id.PopSegment("elasticPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &elasticPool, nil
}

type MssqlVmId struct {
	ResourceGroup string
	Name          string
}

func MssqlVmID(input string) (*MssqlVmId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql VM ID %q: %+v", input, err)
	}

	sqlvm := MssqlVmId{
		ResourceGroup: id.ResourceGroup,
	}

	if sqlvm.Name, err = id.PopSegment("sqlVirtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &sqlvm, nil
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
