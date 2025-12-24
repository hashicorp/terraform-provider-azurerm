# Guide: List Resource

This guide covers how to add a List Resource for an existing resource, using `azurerm_network_profile` as an example. For more information on Lists, see [Resources - List](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/list).

## Prerequisites

Before adding a List Resource, the resource must have Resource Identity implemented. For more information on implementing Resource Identity see [Guide: Resource Identity](guide-resource-identity.md).

## Adding List Resource

1. In the resource, refactor the Read function to have a separate flatten function which will be used in the List Resource.

``` 
func resourceNetworkProfileFlatten(d *pluginsdk.ResourceData, id *networkprofiles.NetworkProfileId, profile *networkprofiles.NetworkProfile) error {
	d.Set("name", id.NetworkProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if profile != nil {
		if props := profile.Properties; props != nil {
			cniConfigs := flattenNetworkProfileContainerNetworkInterface(props.ContainerNetworkInterfaceConfigurations)
			if err := d.Set("container_network_interface", cniConfigs); err != nil {
				return fmt.Errorf("setting `container_network_interface`: %+v", err)
			}

			cniIDs := flattenNetworkProfileContainerNetworkInterfaceIDs(props.ContainerNetworkInterfaces)
			if err := d.Set("container_network_interface_ids", cniIDs); err != nil {
				return fmt.Errorf("setting `container_network_interface_ids`: %+v", err)
			}
		}
		d.Set("location", location.NormalizeNilable(profile.Location))
		if err := tags.FlattenAndSet(d, profile.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}
```

1. Create a new file for the List Resource (for example, `network_profile_resource_list.go`) and scaffold the empty resource:

``` 
type NetworkProfileListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(NetworkProfileListResource)

func (r NetworkProfileListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceNetworkProfile()
}

func (r NetworkProfileListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = `azurerm_network_profile`
}

```

1. Implement the List function. This example uses the subscription ID and resource group as options to list by; other APIs may use different criteria.

```
func (r NetworkProfileListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {

    client := metadata.Client.Network.NetworkProfiles

    // Read the list config data into the model
	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

    // Initialize a list for the results of the API request
	results := make([]networkprofiles.NetworkProfile, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

    // Make the request based on which list parameters have been set in the config
	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureNetworkProfileResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListAllComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureNetworkProfileResourceName), err)
			return
		}

		results = resp.Items
	}

    // Define the function that will push results into the stream 
	stream.Results = func(push func(list.ListResult) bool) {
		for _, profile := range results {
        
            // Initialize a new result object for each resource in the list
			result := request.NewListResult(ctx)
            
            // Set the display name of the item as the resource name
			result.DisplayName = pointer.From(profile.Name)

            // Create a new ResourceData object to hold the state of the resource
			rd := resourceNetworkProfile().Data(&terraform.InstanceState{})
            
            // Set the ID of the resource for the ResourceData object
            id, err := networkprofiles.ParseNetworkProfileID(pointer.From(profile.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Network Profile ID", err)
				return
			}
			rd.SetId(id.ID())

            // Use the resource flatten function to set the attributes into the resource state
			if err := resourceNetworkProfileFlatten(rd, id, &profile); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", azureNetworkProfileResourceName), err)
				return
			}

            // Convert and set the identity and resource state into the result
			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Identity State", err)
				return
			}

			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Identity Data", err)
				return
			}

            // Convert and set the resource state into the result
			tfTypeResourceState, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Resource State", err)
				return
			}

			if err := result.Resource.Set(ctx, *tfTypeResourceState); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Resource Data", err)
				return
			}

            // Send the result to the stream
			if !push(result) {
				return
			}
		}
	}
}

```

1. Register the new List Resource

List Resources are registered within the `registration.go` within each Service Package - and should look something like this:

```
package network

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

var _ sdk.FrameworkServiceRegistration = Registration{}

// ...

// Resources returns a list of List Resources supported by this Service
func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
    return []sdk.FrameworkListWrappedResource{
        NetworkProfileListResource{},
        }
}
```

1. Add Acceptance Tests for this List Resource

Create a new acceptance test file for the List Resource (for example, `network_profile_resource_list_test.go`) and add tests to cover the List Resource functionality. The test should provision any prerequisite resources and multiple resources of the type of List Resource we want to test. 

The test should look something like this:

```
package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccNetworkProfile_list_basic(t *testing.T) {
	r := NetworkProfileResource{}
	listResourceAddress := "azurerm_network_profile.list"

	data := acceptance.BuildTestData(t, "azurerm_network_profile", "test1")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data), // provision multiple resources
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3), // expect at least the 3 we created
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3), // expect exactly the 3 we created in that resource group
				},
			},
		},
	})
}

// provision multiple Network Profile resources for testing
func (r NetworkProfileResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

// Prerequisite Resources ....

resource "azurerm_network_profile" "test1" {
  name                = "acctestnetprofile-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  container_network_interface {
    name = "acctesteth-%[1]d"

    ip_configuration {
      name      = "acctestipconfig-%[1]d"
      subnet_id = azurerm_subnet.test.id
    }
  }
}

resource "azurerm_network_profile" "test2" {
  name                = "acctestnetprofile2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  container_network_interface {
    name = "acctesteth-%[1]d"

    ip_configuration {
      name      = "acctestipconfig-%[1]d"
      subnet_id = azurerm_subnet.test.id
    }
  }
}

resource "azurerm_network_profile" "test3" {
  name                = "acctestnetprofile3-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  container_network_interface {
    name = "acctesteth-%[1]d"

    ip_configuration {
      name      = "acctestipconfig-%[1]d"
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

// define the basic list query for testing
func (r NetworkProfileResource) basicQuery() string {
	return `
list "azurerm_network_profile" "list" {
  provider = azurerm
  config {}
}
`
}

// define the list query for testing by resource group name
func (r NetworkProfileResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_network_profile" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-%[1]d"
  }
}
`, data.RandomInteger)
}

```

1. Add documentation for this List Resource

Documentation should be written manually and added to the `./website/docs/list-resources/` folder.

It should include an example, arguments reference, and look something like this:

````markdown
---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_profile"
description: |-
Lists Network Profile resources.
---

# List resource: azurerm_network_profile

Lists Network Profile resources.

## Example Usage

### List all Network Profiles in the subscription

```hcl
list "azurerm_network_profile" "example" {
  provider = azurerm
  config {}
}
```

### List all Network Profiles in a specific resource group

```hcl
list "azurerm_network_profile" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
````