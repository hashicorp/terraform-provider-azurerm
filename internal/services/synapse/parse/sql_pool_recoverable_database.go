// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SqlPoolRecoverableDatabaseId struct {
	SubscriptionId          string
	ResourceGroup           string
	WorkspaceName           string
	RecoverableDatabaseName string
}

func NewSqlPoolRecoverableDatabaseID(subscriptionId, resourceGroup, workspaceName, recoverableDatabaseName string) SqlPoolRecoverableDatabaseId {
	return SqlPoolRecoverableDatabaseId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		WorkspaceName:           workspaceName,
		RecoverableDatabaseName: recoverableDatabaseName,
	}
}

func (id SqlPoolRecoverableDatabaseId) String() string {
	segments := []string{
		fmt.Sprintf("Recoverable Database Name %q", id.RecoverableDatabaseName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sql Pool Recoverable Database", segmentsStr)
}

func (id SqlPoolRecoverableDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/recoverableDatabases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.RecoverableDatabaseName)
}

// SqlPoolRecoverableDatabaseID parses a SqlPoolRecoverableDatabase ID into an SqlPoolRecoverableDatabaseId struct
func SqlPoolRecoverableDatabaseID(input string) (*SqlPoolRecoverableDatabaseId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SqlPoolRecoverableDatabase ID: %+v", input, err)
	}

	resourceId := SqlPoolRecoverableDatabaseId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.RecoverableDatabaseName, err = id.PopSegment("recoverableDatabases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
