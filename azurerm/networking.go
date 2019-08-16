package azurerm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func extractResourceGroupAndErcName(resourceId string) (resourceGroup string, name string, e error) {
	id, err := azure.ParseAzureResourceID(resourceId)
	if err != nil {
		return "", "", err
	}

	return id.ResourceGroup, id.Path["expressRouteCircuits"], err
}

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

func parseRouteTableName(routeTableId string) (string, error) {
	id, err := azure.ParseAzureResourceID(routeTableId)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to parse Route Table ID '%s': %+v", routeTableId, err)
	}

	return id.Path["routeTables"], nil
}

func retrieveErcByResourceId(resourceId string, meta interface{}) (erc *network.ExpressRouteCircuit, resourceGroup string, e error) {
	ercClient := meta.(*ArmClient).network.ExpressRouteCircuitsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, name, err := extractResourceGroupAndErcName(resourceId)
	if err != nil {
		return nil, "", fmt.Errorf("Error Parsing Azure Resource ID -: %+v", err)
	}

	resp, err := ercClient.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("Error making Read request on Express Route Circuit %s: %+v", name, err)
	}

	return &resp, resGroup, nil
}
