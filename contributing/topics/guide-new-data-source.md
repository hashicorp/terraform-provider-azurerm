# Guide: New Data Source

This guide covers adding a new Data Source to a Service Package, see [adding a New Service Package](guide-new-service-package.md) if the Service Package doesn't exist yet.

### Related Topics

* [Acceptance Testing](reference-acceptance-testing.md)
* [Our Recommendations for opening a Pull Request](guide-opening-a-pr.md)

### Stages

At this point in time the AzureRM Provider supports both Typed and Untyped Data Sources - more information can be found [in the High Level Overview](high-level-overview.md).

This guide covers adding a new Typed Data Source, which makes use of the Typed SDK within this repository and requires the following steps:

1. Ensure all the dependencies are installed (see [Building the Provider](building-the-provider.md)).
2. Add an SDK Client (if required).
3. Define the Resource ID.
4. Scaffold an empty/new Data Source.
5. Register the new Data Source.
6. Add Acceptance Test(s) for this Data Source.
7. Run the Acceptance Test(s).
8. Add Documentation for this Data Source.
9. Send the Pull Request.

We'll go through each of those steps in turn, presuming that we're creating a Data Source for a Resource Group.

### Step 1: Ensure the Tools are installed

See [Building the Provider](building-the-provider.md).

### Step 2: Add an SDK Client (if required)

If you're creating a new Data Source for a Resource that's already created by Terraform, the SDK Client you need to use is likely already supported (and so you can skip this section).

However if the SDK Client you need to use isn't already configured in the Provider, we'll cover how to add and configure the SDK Client.

Determining which SDK Client you should be using is a little complicated unfortunately, in this case the SDK Client we want to use is: `github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources`.

The Client for the Service Package can be found in `./internal/services/{name}/client/client.go` - and we can add an instance of the SDK Client we want to use (here `resources.GroupsClient`) and configure it (adding credentials etc): 

```go
package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/resources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GroupsClient *resources.GroupsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	groupsClient, err := resources.NewResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Resources Client: %+v", err)
    }
	o.Configure(groupsClient.Client, o.Authorizer.ResourceManager)
	
	// ...
	
	return &Client{
		GroupsClient: groupsClient,
	}
}
```

Things worth noting here:

- The call to `o.Configure` configures the authorization token which should be used for this SDK Client - in most cases `ResourceManager` is the authorizer you want to use.

At this point, this SDK Client should be usable within the Data Sources via:

```go
client := metadata.Client.{ServicePackage}.{ClientField}
```

For example, in this case:

```go
client := metadata.Client.Resource.GroupsClient
```

### Step 3: Scaffold an empty/new Data Source

Since we're creating a Data Source for a Resource Group, which is a part of the Resources API - we'll want to create an empty Go file within the Service Package for Resources, which is located at `./internal/services/resources`.

In this case, this would be a file called `resource_group_example_data_source.go`, which we'll start out with the following:

> **Note:** We'd normally name this file `resource_group_data_source.go` - but there's an existing Data Source for Resource Groups, so we're appending `example` to the name throughout this guide. 

```go
package resources

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.DataSource = ResourceGroupExampleDataSource{}

type ResourceGroupExampleDataSource struct {}
```

> **Note:** Your editor may show a suggestion to implement the methods defined in `sdk.DataSource` for the `ResourceGroupExampleDataSource` struct - we'd recommend holding off the first time around to explain each of the methods.

In this case the interface `sdk.DataSource` defines all of the methods required for a Data Source which the newly created struct for the Resource Group Data Source need to implement, which are:

```go
type DataSource interface {
    Arguments() map[string]*schema.Schema
    Attributes() map[string]*schema.Schema
    ModelObject() interface{}
    ResourceType() string
	Read() ResourceFunc
}
```

To go through these in turn:

* `Arguments` returns a list of schema fields which are user-specifiable - either Required or Optional.
* `Attributes` returns a list of schema fields which are Computed (read-only).
* `ModelObject` returns a reference to a Go struct which is used as the Model for this Data Source (this can also return `nil` if there's no model).
* `ResourceType` returns the name of this resource within the Provider (for example `azurerm_resource_group_example`).
* `Read` returns a function defining both the Timeout and the Read function (which retrieves information from the Azure API) for this Data Source.

```go
func (ResourceGroupExampleDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (ResourceGroupExampleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (ResourceGroupExampleDataSource) ModelObject() interface{} {
	return nil
}

func (ResourceGroupExampleDataSource) ResourceType() string {
	return "azurerm_resource_group_example"
}
```

> In this case we're using the resource type `azurerm_resource_group_example` as [an existing Data Source for `azurerm_resource_group` exists](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/resource_group) and the names need to be unique.

These functions define a Data Source called `azurerm_resource_group_example`, which has one Required argument called `name` and two Computed arguments called `location` and `tags`. We'll come back to `ModelObject` later.

---

Next up, let's implement the Read function - which retrieves the information about the Resource Group from Azure:

```go
func (ResourceGroupExampleDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 5 minutes may initially seem excessive, we set this as a default to account for rate
		// limiting - but having this here means that users can override this in their config as necessary
		Timeout: 5 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state 
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.Resource.GroupsClient
            
			// retrieve the Name for this Resource Group from the Terraform Config
			// and then create a Resource ID for this Resource Group
			// using the Subscription ID & name
            subscriptionId := metadata.Client.Account.SubscriptionId
            name := metadata.ResourceData.Get("name").(string)
            id := resources.NewResourceGroupExampleID(subscriptionId, name)
			
			// then retrieve the Resource Group by it's ID
            resp, err := client.Get(ctx, id)
            if err != nil {
				// if the Resource Group doesn't exist (e.g. we get a 404 Not Found)
				// since this is a Data Source we must return an error if it's Not Found
                if response.WasNotFound(resp.HttpResponse) {
                    return fmt.Errorf("%s was not found", id)
                }
				
                // otherwise it's a genuine error (auth/api error etc) so raise it
				// there should be enough context for the user to interpret the error
				// or raise a bug report if there's something we should handle
                return fmt.Errorf("retrieving %s: %+v", id, err)
            }
			
			// now we know the Resource Group exists, set the Resource ID for this Data Source
			// this means that Terraform will track this as existing
            metadata.SetID(id)
			
			// at this point we can set information about this Resource Group into the State
			// whilst traditionally we would do this via `metadata.ResourceData.Set("foo", "somevalue")
			// the Location and Tags fields are a little different - and we have a couple of normalization
			// functions for these.
			
			// whilst this may seem like a weird thing to call out in an example, because these two fields
			// are present on the majority of resources, we hope it explains why they're a little different
			 
			// in this case the Location can be returned in various different forms, for example
			// "West Europe", "WestEurope" or "westeurope" - as such we normalize these into a
			// lower-cased singular word with no spaces (e.g. "westeurope") so this is consistent
			// for users
			if model := resp.Model; model != nil {
			    metadata.ResourceData.Set("location", location.NormalizeNilable(model.Location))
				
				props := model.Properties; props != nil {
					// If the data source exposes additional properties that live within the Properties
					// model of the response they would be set into state here.
                }
                // (as above) Tags are a little different, so we have a dedicated helper function
                // to flatten these consistently across the Provider
                return tags.FlattenAndSet(metadata.ResourceData, model.Tags)
            }   
			return nil
		},
	}
}
```

---

At this point the finished Data Source should look like (including imports):

```go
package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ResourceGroupExampleDataSource struct{}

func (d ResourceGroupExampleDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (d ResourceGroupExampleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"tags": commonschema.TagsDataSource(),

	}
}

func (d ResourceGroupExampleDataSource) ModelObject() interface{} {
	return nil
}

func (d ResourceGroupExampleDataSource) ResourceType() string {
	return "azurerm_resource_group_example"
}

func (d ResourceGroupExampleDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			
			id := resources.NewResourceGroupExampleID(subscriptionId, metadata.ResourceData.Get("name").(string))

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				metadata.ResourceData.Set("location", location.NormalizeNilable(model.Location))
				return tags.FlattenAndSet(metadata.ResourceData, model.Tags)
			}
			return nil
		},
	}
}
```

At this point in time this Data Source is now code-complete - there's an optional extension to make this cleaner by using a Typed Model, however this isn't necessary.

### Step 4: Register the new Data Source

Data Sources are registered within the `registration.go` within each Service Package - and should look something like this:

```go
package resource

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

// ...

// DataSources returns a list of Data Sources supported by this Service
func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}
```

---

> **Note:** It's possible that the Service Registration (above) doesn't currently support Typed Resources, in which case you may need to add the following:

```go
var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct {
}

func (Registration) Name() string {
	return "Some Service"
}

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}

func (Registration) WebsiteCategories() []string {
	return []string{
		"Some Service",
	}
}
```

> In this case you'll also need to add a line to register this Service Registration [in the list of Typed Service Registrations](https://github.com/hashicorp/terraform-provider-azurerm/blob/bd7c755b789fa131778ef93824cf3bae5caccf56/internal/provider/services.go#L109).

---

To register the Data Source we need to add an instance of the struct used for the Data Source to the list of Data Sources, for example:

```go
// DataSources returns a list of Data Sources supported by this Service
func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ResourceGroupExampleDataSource{},	
	}
}
```

At this point the Data Source is registered, as when the Azure Provider builds up a list of supported Data Sources during initialization, it parses each of the Service Registrations to put together a definitive list of the Data Sources that we support.

This means that if you [Build the Provider](building-the-provider.md), at this point you should be able to apply the following Data Source:

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_resource_group_example" "test" {
  name = "some-pre-existing-resource-group" # presuming this resource group exists ;)
}

output "location" {
  value = data.azurerm_resource_group_example.test.location
}
```

### Step 5: Add Acceptance Test(s) for this Data Source

We're going to test the Data Source that we've just built by dynamically provisioning a Resource Group using the Azure Provider, then asserting that we can look up that Resource Group using the new `azurerm_resource_group_example` Data Source.

In Go tests are expected to be in a file name in the format `{original_file_name}_test.go` - in our case that'd be `resource_group_example_data_source_test.go`, into which we'll want to add: 

```go
package resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ResourceGroupExampleDataSource struct{}

func TestAccResourceGroupExampleDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resource_group_example", "test")
	r := ResourceGroupExampleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
	})
}

func (ResourceGroupExampleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"

  tags = {
    env = "test"
  }
}

data "azurerm_resource_group_example" "test" {
  name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
```

There's a more detailed breakdown of how this works [in the Acceptance Testing reference](reference-acceptance-testing.md) - but to summarize what's going on here:

1. Test Terraform Configurations are defined as methods on the struct `ResourceGroupExampleDataSource` so that they're easily accessible (this helps to avoid them being unintentionally used in other resources).
2. The `acceptance.TestData` object contains a number of helpers, including both random integers, strings and the Azure Locations where resources should be provisioned - which are used to ensure when tests are run in parallel that we provision unique resources for testing purposes.
3. We're asserting on the Computed (e.g. read-only) fields returned from the Resource - we don't check the user-specified fields (`name` in this case) as if it's missing, the test will fail to find the Resource Group.
4. We append `_test` to the Go package name (e.g. `resource_test`) since we need to be able to access both the `resource` package and the `acceptance` package (which is a circular reference, otherwise).

At this point we should be able to run this test.

### Step 6: Run the Acceptance Test(s)

Detailed [instructions on Running the Tests can be found in this guide](running-the-tests.md) - when a Service Principal is configured you can run the test above using:

```sh
make acctests SERVICE='resource' TESTARGS='-run=TestAccResourceGroupExampleDataSource_basic' TESTTIMEOUT='60m'
```

Which should output:

```sh
==> Checking that code complies with gofmt requirements...
==> Checking that Custom Timeouts are used...
==> Checking that acceptance test packages are used...
TF_ACC=1 go test -v ./internal/services/resource -run=TestAccResourceGroupExampleDataSource_basic -timeout 60m -ldflags="-X=github.com/hashicorp/terraform-provider-azurerm/version.ProviderVersion=acc"
=== RUN   TestAccResourceGroupExampleDataSource_basic
=== PAUSE TestAccResourceGroupExampleDataSource_basic
=== CONT  TestAccResourceGroupExampleDataSource_basic
--- PASS: TestAccResourceGroupExampleDataSource_basic (88.15s)
PASS
ok  	github.com/hashicorp/terraform-provider-azurerm/internal/services/resource	88.735s
```

### Step 7: Add Documentation for this Data Source

At this point in time documentation for each Data Source (and Resource) is written manually, located within the `./website` folder - in this case this will be located at `./website/docs/d/resource_group_example.html.markdown`.

There is a tool within the repository to help scaffold the documentation for a Data Source - the documentation for this Data Source can be scaffolded via the following command:

```sh
$ make scaffold-website BRAND_NAME="Resource Group Example" RESOURCE_NAME="azurerm_resource_group_example" RESOURCE_TYPE="data"
```

The documentation should look something like below - containing both an example usage and the required, optional and computed fields:

> **Note:** In the example below you'll need to replace each `[]` with a backtick "`" - as otherwise this gets rendered incorrectly, unfortunately.

```markdown
---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_resource_group_example"
description: |-
  Gets information about an existing Resource Group.
---

# Data Source: azurerm_resource_group_example

Use this data source to access information about an existing Resource Group.

## Example Usage

[][][]hcl
data "azurerm_resource_group_example" "example" {
  name = "existing"
}

output "id" {
  value = data.azurerm_resource_group_example.example.id
}
[][][]

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name of this Resource Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Group.

* `location` - The Azure Region where the Resource Group exists.

* `tags` - A mapping of tags assigned to the Resource Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.
```

> **Note:** In the example above you'll need to replace each `[]` with a backtick "`" - as otherwise this gets rendered incorrectly, unfortunately.

### Step 8: Send the Pull Request

See [our recommendations for opening a Pull Request](guide-opening-a-pr.md).
