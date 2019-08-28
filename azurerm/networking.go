package azurerm

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func evaluateSchemaValidateFunc(i interface{}, k string, validateFunc schema.SchemaValidateFunc) (bool, error) { // nolint: unparam
	_, errors := validateFunc(i, k)

	errorStrings := []string{}
	for _, e := range errors {
		errorStrings = append(errorStrings, e.Error())
	}

	if len(errors) > 0 {
		return false, fmt.Errorf(strings.Join(errorStrings, "\n"))
	}

	return true, nil
}

func parseNetworkSecurityGroupName(networkSecurityGroupId string) (string, error) {
	id, err := azure.ParseAzureResourceID(networkSecurityGroupId)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to Parse Network Security Group ID '%s': %+v", networkSecurityGroupId, err)
	}

	return id.Path["networkSecurityGroups"], nil
}
