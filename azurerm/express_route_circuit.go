package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/errwrap"
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
		return nil, "", errwrap.Wrapf("Error Parsing Azure Resource ID - {{err}}", err)
	}

	resp, err := ercClient.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, "", nil
		}
		return nil, "", errwrap.Wrapf(fmt.Sprintf("Error making Read request on Express Route Circuit %s: {{err}}", name), err)
	}

	return &resp, resGroup, nil
}
