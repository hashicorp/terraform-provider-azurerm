// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

// @tombuildsstuff: note this can't be automated today due to the comma in the Resource ID

type RestorableDroppedDatabaseId struct {
	Name          string
	MsSqlServer   string
	ResourceGroup string
	RestoreName   string
}

func RestorableDroppedDatabaseID(input string) (*RestorableDroppedDatabaseId, error) {
	inputList := strings.Split(input, ",")

	if len(inputList) != 2 {
		return nil, fmt.Errorf("[ERROR] Unable to parse Microsoft Sql Restorable DB ID %q, please refer to '/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/sqlServer1/restorableDroppedDatabases/sqlDB1,000000000000000000'", input)
	}

	restorableDBId := RestorableDroppedDatabaseId{
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
