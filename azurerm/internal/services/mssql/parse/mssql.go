package parse

import (
	"fmt"
	"strings"

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

type MsSqlRestorableDBId struct {
	Name          string
	MsSqlServer   string
	ResourceGroup string
	RestoreName   string
}

type MsSqlRecoverableDBId struct {
	Name          string
	MsSqlServer   string
	ResourceGroup string
}

type MsSqlSyncGroupId struct {
	ResourceGroup string
	ServerName    string
	DatabaseName  string
	Name          string
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

func MssqlRestorableDBID(input string) (*MsSqlRestorableDBId, error) {
	inputList := strings.Split(input, ",")

	if len(inputList) != 2 {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql Restorable DB ID %q, please refer to '/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/restorableDroppedDatabases/sqlDB1,000000000000000000'", input)
	}

	restorableDBId := MsSqlRestorableDBId{
		RestoreName: inputList[1],
	}

	id, err := azure.ParseAzureResourceID(inputList[0])
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql Restorable DB ID %q: %+v", input, err)
	}

	restorableDBId.ResourceGroup = id.ResourceGroup

	if restorableDBId.MsSqlServer, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if restorableDBId.Name, err = id.PopSegment("restorableDroppedDatabases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(inputList[0]); err != nil {
		return nil, err
	}

	return &restorableDBId, nil
}

func MssqlRecoverableDBID(input string) (*MsSqlRecoverableDBId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql Recoverable DB ID %q: %+v", input, err)
	}

	recoverableDBId := MsSqlRecoverableDBId{
		ResourceGroup: id.ResourceGroup,
	}

	if recoverableDBId.MsSqlServer, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}

	if recoverableDBId.Name, err = id.PopSegment("recoverabledatabases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &recoverableDBId, nil
}

func MsSqlSyncGroupID(input string) (*MsSqlSyncGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing SQL Sync Group ID %q: %+v", input, err)
	}

	syncGroupId := MsSqlSyncGroupId{
		ResourceGroup: id.ResourceGroup,
	}
	if syncGroupId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if syncGroupId.DatabaseName, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}
	if syncGroupId.Name, err = id.PopSegment("syncGroups"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &syncGroupId, nil
}