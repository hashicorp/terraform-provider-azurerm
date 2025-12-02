# Guide: Resource IDs

Resource IDs are an essential component of all resources and data sources within the provider and are required in order to interact with the corresponding API and to be tracked in state correctly by Terraform.

Due to their fundamental importance in the provider, resource IDs should be handled through the use of various helper functions that are available.

## Resource ID Parsers and Validators in `hashicorp/go-azure-sdk`

The SDK `hashicorp/go-azure-sdk` used in the provider contains the necessary parsing and validation functions required for a resource ID and will exist in the resource's package.

As an example the function to parse a Machine Learning Workspace Resource ID will be accessible by importing the workspace resource package into the provider:

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"

func (r MachineLearningWorkspace) Create() sdk.ResourceFunc {
	...
	id := workspaces.NewWorkspaceID(subscriptionId, resourceGroupName, workspaceName)
	...
}

func (r MachineLearningWorkspace) Read() sdk.ResourceFunc {
    ...
    id := workspaces.ParseWorkspaceID(subscriptionId, resourceGroupName, workspaceName)
	...
}
```

## Resource ID Parsers and Validators from `hashicorp/go-azure-helpers/resourcemanager/commonids`

Some resource types that are referenced across multiple services, will have their parser and validation functions defined in `hashicorp/go-azure-helpers` `commonids` package.

This is done to avoid having to convert the same ID between different types.

```go
import `github.com/hashicorp/go-azure-helpers/resourcemanager/commonids`

func (r AppServicePlan) Create() sdk.ResourceFunc {
    ...
    id := commonids.NewAppServicePlanID(subscriptionId, resourceGroupName, workspaceName)
	...
}

func (r AppServicePlan) Read() sdk.ResourceFunc {
    ...
    id := commonids.ParseAppServicePlanID(subscriptionId, resourceGroupName, workspaceName)
	...
}
```

## Composite Resource IDs

Resource IDs that consist of a scope ID and a resource ID separated by a `|` e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/publicIPAddresses/myPublicIpAddress1` should be handled using the Composite Resource ID functions in the `hashicorp/go-azure-helpers` `commonids` package.

```go
import `github.com/hashicorp/go-azure-helpers/resourcemanager/commonids`

func (r NatGatewayPublicIpAssociation) Create() sdk.ResourceFunc {
    ...

	publicIpAddressId, err := commonids.ParsePublicIPAddressID(d.Get("public_ip_address_id").(string))
    if err != nil {
		return err
    }
    
    natGatewayId, err := natgateways.ParseNatGatewayID(d.Get("nat_gateway_id").(string))
    if err != nil {
        return err
    }

    id := commonids.NewCompositeResourceID(natGatewayId, publicIpAddressId)
	...
}

func (r NatGatewayPublicIpAssociation) Read() sdk.ResourceFunc {
    ...
	id, err := commonids.ParseCompositeResourceID(d.Id(), &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
    if err != nil {
        return err
    }
	...
    d.Set("nat_gateway_id", id.First.ID())
    d.Set("public_ip_address_id", id.Second.ID())
	...
}
```

## Generated Resource ID Parsers and Validators (legacy)

Prior to generating the parser and validation functions within the SDK, we generated these functions in the provider with [this automation](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/internal/tools/generator-resource-id) which generates the functions for all IDs defined in `resourceids.go`.

An example of this is shown below:

```go
package resource

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroupExample -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1
```

In this case, you need to specify the `name` of the Resource (in this case `ResourceGroupExample`) and the `id` which is an example of this Resource ID (in this case `/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1`).

Running `make generate` - will output the following files:

* `./internal/service/resource/parse/resource_group_example.go` - contains the Resource ID Struct, Formatter and Parser.
* `./internal/service/resource/parse/resource_group_example_test.go` - contains tests for those ^.
* `./internal/service/resource/validate/resource_group_example_id.go` - contains Terraform validation functions for the Resource ID.

> **Note:** This is an outdated way of handling resource IDs in the provider and is being phased out. This method should only be used in exceptional cases.

