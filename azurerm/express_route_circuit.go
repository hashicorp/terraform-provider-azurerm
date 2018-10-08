package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-08-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func extractResourceGroupAndErcName(resourceId string) (resourceGroup string, name string, err error) {
	id, err := parseAzureResourceID(resourceId)

	if err != nil {
		return "", "", err
	}
	resourceGroup = id.ResourceGroup
	name = id.Path["expressRouteCircuits"]

	return
}

func retrieveErcByResourceId(resourceId string, meta interface{}) (erc *network.ExpressRouteCircuit, resourceGroup string, e error) {
	ercClient := meta.(*ArmClient).expressRouteCircuitClient
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
