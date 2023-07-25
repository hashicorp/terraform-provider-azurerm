// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SqlFunctionId struct {
	SubscriptionId          string
	ResourceGroup           string
	DatabaseAccountName     string
	SqlDatabaseName         string
	ContainerName           string
	UserDefinedFunctionName string
}

func NewSqlFunctionID(subscriptionId, resourceGroup, databaseAccountName, sqlDatabaseName, containerName, userDefinedFunctionName string) SqlFunctionId {
	return SqlFunctionId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		DatabaseAccountName:     databaseAccountName,
		SqlDatabaseName:         sqlDatabaseName,
		ContainerName:           containerName,
		UserDefinedFunctionName: userDefinedFunctionName,
	}
}

func (id SqlFunctionId) String() string {
	segments := []string{
		fmt.Sprintf("User Defined Function Name %q", id.UserDefinedFunctionName),
		fmt.Sprintf("Container Name %q", id.ContainerName),
		fmt.Sprintf("Sql Database Name %q", id.SqlDatabaseName),
		fmt.Sprintf("Database Account Name %q", id.DatabaseAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sql Function", segmentsStr)
}

func (id SqlFunctionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s/containers/%s/userDefinedFunctions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.UserDefinedFunctionName)
}

// SqlFunctionID parses a SqlFunction ID into an SqlFunctionId struct
func SqlFunctionID(input string) (*SqlFunctionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SqlFunction ID: %+v", input, err)
	}

	resourceId := SqlFunctionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}
	if resourceId.SqlDatabaseName, err = id.PopSegment("sqlDatabases"); err != nil {
		return nil, err
	}
	if resourceId.ContainerName, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}
	if resourceId.UserDefinedFunctionName, err = id.PopSegment("userDefinedFunctions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
