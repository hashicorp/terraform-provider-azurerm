# Acceptance Testing

Acceptance tests are an essential part of the provider - they provide confidence in the functionality and consistency of resources and data sources as they are introduced and over time.

Whilst we can't test every use-case or permutation of fields - each data source/resource gets a common set of tests to ensure the core use-cases are covered.

As a general rule, the more complex the resource the more tests there are - for example AKS, App Service and Virtual Machines all have a large number of end-to-end tests.

### Considerations

> **Note:** Acceptance Tests provision real resources within Azure - which may have an associated charge for each resource.

* When selecting SKUs for testing, **pick the lowest or cheapest SKU** which covers the test - unless there's good reason to otherwise (e.g. some configurations can provision more quickly using one SKU over another).

* Always put the resource being tested **at the end of each configuration**, especially if a test requires multiple resource or data declarations. This makes it easy to find the resource being tested, especially in large configurations. Example: note how the `azurerm_virtual_machine` is positioned at the end of each configuration [virtual_machine_resource_test.go](../../internal/services/legacy/virtual_machine_resource_test.go)

* Use a fixed password that meets the service's rules rather than generating one, unless the test is about generated passwords.

### Running the Tests

See [Running the Tests](running-the-tests.md).

### PreCheck Helpers

Some tests need something not every subscription has, such as a real DNS zone or an existing Key Vault certificate. Without a check up front, the test starts creating resources it can never finish.

#### Global pre-check

`acceptance.PreCheck(t)` already checks credentials and test locations. It fails outright, which is right for something every test needs.

#### Resource-specific pre-check

Give the test struct a `preCheck` method that skips the test when the environment variable for it is unset, and call it at the start of each test that needs it. Skip rather than fail, so contributors without that infrastructure can still run everything else while CI supplies it and runs the full suite.

```go
type ExampleResource struct{}

func TestAccExampleResource_basic(t *testing.T) {
        data := acceptance.BuildTestData(t, "azurerm_example_resource", "test")
        r := ExampleResource{}
        r.preCheck(t)

        data.ResourceTest(t, r, []acceptance.TestStep{
                {
                        Config: r.basic(data),
                        Check: acceptance.ComposeTestCheckFunc(
                                check.That(data.ResourceName).ExistsInAzure(r),
                        ),
                },
                data.ImportStep(),
        })
}

func (ExampleResource) preCheck(t *testing.T) {
        if os.Getenv("ARM_TEST_SOME_PREREQ") == "" {
                t.Skip("Skipping as ARM_TEST_SOME_PREREQ is not set")
        }
}
```

#### Where to put `preCheck`

Next to the tests that call it, after the `TestAcc...` functions and before `Exists`.

### Test Package

While tests reside in the same folder as resource and data source .go files, they need to be in a separate test package to prevent circular references. i.e. for the file `./internal/services/aadb2c/aadb2c_directory_data_source_test.go` the package should be:

```go
package aadb2c_test

import ...
```

This is checked by `make test` during CI.

### Import Step

During acceptance tests it is important to validate that the resource in Azure matches what Terraform expects and has saved into state. This can be done by adding a `data.ImportStep()` after every step. This will import the resource into Terraform and compare that the Terraform state matches the Azure Resource.

As some properties (such as sensitive data like passwords) are not returned from Azure you can ignore these properties by passing them into the import step: `data.ImportStep("password", "database_primary_key")`.

### Naming

Test names should follow the convention `TestAcc` + `ResourceName` + `_` + `test` -> `TestAccExampleResource_basic`, or to group tests:

```go
func TestAccExampleResource_category_test1(t *testing.T) { ... }
func TestAccExampleResource_category_test2(t *testing.T) { ... }
```

---

## Acceptance Tests

The Acceptance Tests for both Data Sources and Resources within this Provider use a Go struct for each test, in the form `{Name}{DataSource|Resource}`, for example:

```go
// for a data source named Example:
type ExampleDataSource struct {}

// for a resource named Example:
type ExampleResource struct {}
```

They are differentiated from the implementation's struct by their package, which is the same as the implementation's but with a `_test` suffix. This allows the test configurations to be scoped (and not used unintentionally across different resources), for example a Resource may look like this:

```go
package example_test

type ExampleResource struct {}

func (ExampleResource) basic(data acceptance.TestData) string {
return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_example_resource" "example" {
  name             = "my_example_resource"
  location         = "%s"
  example_property = "bar"
}

`, data.Locations.Primary)
}
```

This allows the Acceptance Test for each Data Source/Resource to reference that struct and obtain the associated Terraform Configuration as a part of the test e.g.:

```go
func TestAccExampleResource_basic(t *testing.T) {
        data := acceptance.BuildTestData(t, "azurerm_example_resource", "test")
        r := ExampleResource{}

        data.ResourceTest(t, r, []acceptance.TestStep{
                {
                        Config: r.basic(data),
                        Check: acceptance.ComposeTestCheckFunc(
                                check.That(data.ResourceName).ExistsInAzure(r),
                        ),
                },
                data.ImportStep(),
        })
}
```

### Which Tests are Required?

At a minimum, a Data Source requires:

* A `basic` test ([Example](#example---data-source---basic)) - this tests the minimum fields (e.g. all Required fields) for this Data Source.

However, more complex Data Sources can warrant additional acceptance tests - consideration should be given during the development of each Data Source to what's important to be tested.

---

At a minimum, a Resource requires:

* A `basic` test ([Example](#example---resource---basic)) - this tests the minimum fields (e.g. all Required fields) for this Resource.

* A `requiresImport` test ([Example](#example---resource---requires-import)) - creates the resource twice and expects the second attempt to fail, using `RequiresImportErrorStep`.

* A `complete` test ([Example](#example---resource---complete)) - this tests all possible fields (e.g. all Required/Optional fields) for this Resource.

* A `update` test ([Example](#example---resource---update)) - This test exercises a change of values for any properties that can be updated by executing consecutive configurations to change a resource in a predictable manner. Properties which are `ForceNew` should not be tested in this way.

However, more complex Resource generally warrant additional acceptance tests - consideration should be given during the development of each Resource to what's important to be tested.

### Example - Data Source - Basic

Reuse the resource's `complete` configuration, so every attribute the data source exposes has a value to assert, and check those values. There is no import step, since a data source has nothing to import.

When one config helper is only being threaded into `fmt.Sprintf` once, pass it directly as the argument instead of assigning a temporary variable first.

```go
func TestAccExampleDataSource_basic(t *testing.T) {
        data := acceptance.BuildTestData(t, "data.azurerm_example_resource", "test")
        r := ExampleDataSource{}

        data.DataSourceTest(t, []acceptance.TestStep{
                {
                        Config: r.basic(data),
                        Check: acceptance.ComposeTestCheckFunc(
                            check.That(data.ResourceName).Key("example_property").HasValue("bar"),
                            check.That(data.ResourceName).Key("example_optional_bool").HasValue("false"),
                            check.That(data.ResourceName).Key("example_optional_string").HasValue("foo"),
                        ),
                },
        })
}

func (ExampleDataSource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
%[1]s

data "azurerm_example_resource" "test" {
  name = azurerm_example_resource.test.name
}
`, ExampleResource{}.complete(data))
}
```

---

### Example - Resource - Basic

This test provisions the resource using the minimum configuration possible (e.g. only the `Required` fields), which is intended to test the happy path (of creating, reading and then destroying a resource).

As we're testing the Resource, we make use of an `ImportStep` as a part of the Acceptance Test to ensure that each of the fields specified as a part of the Terraform Configuration are set into the state.

```go
func TestAccExampleResource_basic(t *testing.T) {
        data := acceptance.BuildTestData(t, "azurerm_example_resource", "test")
        r := ExampleResource{}

        data.ResourceTest(t, r, []acceptance.TestStep{
                {
                        Config: r.basic(data),
                        Check: acceptance.ComposeTestCheckFunc(
                                check.That(data.ResourceName).ExistsInAzure(r),
                        ),
                },
                data.ImportStep(),
        })
}

func (ExampleResource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_example_resource" "example" {
  name             = "my_example_resource"
  location         = "%s"
  example_property = "bar"
}
`, data.Locations.Primary)
}
```

### Example - Resource - Complete

This test provisions the resource using the maximum configuration possible (e.g. all `Required` and `Optional` fields which can be set together), which is intended to test the more complex scenario for this resource.

As we're testing the Resource, we make use of an `ImportStep` as a part of the Acceptance Test to ensure that each of the fields specified as a part of the Terraform Configuration are set into the state.

```go
func TestAccExampleResource_complete(t *testing.T) {
        data := acceptance.BuildTestData(t, "azurerm_example_resource", "test")
        r := ExampleResource{}

        data.ResourceTest(t, r, []acceptance.TestStep{
                {
                        Config: r.complete(data),
                        Check: acceptance.ComposeTestCheckFunc(
                                check.That(data.ResourceName).ExistsInAzure(r),
                        ),
                },
                data.ImportStep(),
        })
}

func (ExampleResource) complete(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_example_resource" "example" {
  name             = "my_example_resource"
  location         = "%s"
  example_property = "bar"

  example_optional_bool   = false
  example_optional_string = "foo"

  tags = {
    "Hello" = "World"
  }
}

`, data.Locations.Primary)
}
```

### Example - Resource - Requires Import

Azure would silently overwrite an existing resource, so the create function checks for one first. This test creates the resource with the `basic` config, then tries to create it again with the same name and expects the error.

```go
func TestAccExampleResource_requiresImport(t *testing.T) {
        data := acceptance.BuildTestData(t, "azurerm_example_resource", "test")
        r := ExampleResource{}

        data.ResourceTest(t, r, []acceptance.TestStep{
                {
                        Config: r.basic(data),
                        Check: acceptance.ComposeTestCheckFunc(
                                check.That(data.ResourceName).ExistsInAzure(r),
                        ),
                },
                data.RequiresImportErrorStep(r.requiresImport),
        })
}

func (r ExampleResource) requiresImport(data acceptance.TestData) string {
        return fmt.Sprintf(`
%[1]s

resource "azurerm_example_resource" "import" {
  name             = azurerm_example_resource.example.name
  location         = azurerm_example_resource.example.location
  example_property = azurerm_example_resource.example.example_property
}
`, r.basic(data))
}
```

### Example - Resource - Update

This test is used to confirm that the `Update` function of the Resource works - as such only `Required`/`Optional` fields which are not `ForceNew` can be updated.

The bare-minimum example for this is provisioning the `basic` configuration and then updating it using the `complete` test configuration above, for example:

```go
func TestAccExampleResource_update(t *testing.T) {
        data := acceptance.BuildTestData(t, "azurerm_example_resource", "test")
        r := ExampleResource{}

        data.ResourceTest(t, r, []acceptance.TestStep{
            {   // first provision the resource
                Config: r.basic(data),
                Check: acceptance.ComposeTestCheckFunc(
                    check.That(data.ResourceName).ExistsInAzure(r),
                ),
            },
            data.ImportStep(),
            {   // then perform the update
                Config: r.complete(data),
                Check: acceptance.ComposeTestCheckFunc(
                    check.That(data.ResourceName).ExistsInAzure(r),
                ),
            },
            data.ImportStep(),
        })
}
```

However, this doesn't necessarily cover all use-cases for this resource - or may be too broad depending on the resource, as such it's also common to have tests covering a subset of the fields, for example:

> **Note:** This is a simplified example for testing purposes, we'd generally recommend a test covering a related subset of the resource (e.g. enabling/disabling a block within the resource), rather than a single field - but it depends on the resource.

```go
func TestAccExampleResource_someSetting(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_example_resource", "test")
    r := ExampleResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {   // first provision the resource
            Config: r.someSetting(data, true),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
        {   // then perform the update to disable this setting
            Config: r.someSetting(data, false),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
        {   // finally, check we can re-enable this once it's been disabled
            Config: r.someSetting(data, true),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

func (ExampleResource) someSetting(data acceptance.TestData, enabled bool) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_example_resource" "example" {
  name                 = "my_example_resource"
  location             = "%[1]s"
  some_setting_enabled = %[2]t
}
`, data.Locations.Primary, enabled)
}
```
