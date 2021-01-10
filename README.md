# Terraform Provider for Azure (Resource Manager)

Version 2.x of the AzureRM Provider requires Terraform 0.12.x and later.

* [Terraform Website](https://www.terraform.io)
* [AzureRM Provider Documentation](https://www.terraform.io/docs/providers/azurerm/index.html)
* [AzureRM Provider Usage Examples](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples)
* [Slack Workspace for Contributors](https://terraform-azure.slack.com) ([Request Invite](https://join.slack.com/t/terraform-azure/shared_invite/enQtNDMzNjQ5NzcxMDc3LWNiY2ZhNThhNDgzNmY0MTM0N2MwZjE4ZGU0MjcxYjUyMzRmN2E5NjZhZmQ0ZTA1OTExMGNjYzA4ZDkwZDYxNDE))

## Usage Example

```
# Configure the Microsoft Azure Provider
provider "azurerm" {
  # We recommend pinning to the specific version of the Azure Provider you're using
  # since new versions are released frequently
  version = "=2.40.0"

  features {}

  # More information on the authentication methods supported by
  # the AzureRM Provider can be found here:
  # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs

  # subscription_id = "..."
  # client_id       = "..."
  # client_secret   = "..."
  # tenant_id       = "..."
}

# Create a resource group
resource "azurerm_resource_group" "example" {
  name     = "production-resources"
  location = "West US"
}

# Create a virtual network in the production-resources resource group
resource "azurerm_virtual_network" "test" {
  name                = "production-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}
```

Further [usage documentation is available on the Terraform website](https://www.terraform.io/docs/providers/azurerm/index.html).

## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.12.x +
* [Go](https://golang.org/doc/install) version 1.15.x (to build the provider plugin)

### On Windows

If you're on Windows you'll also need:
* [Git Bash for Windows](https://git-scm.com/download/win)
* [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.*

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".*

Or install via [Chocolatey](https://chocolatey.org/install) (`Git Bash for Windows` must be installed per steps above)
```powershell
choco install make golang terraform -y
refreshenv
```

You must run  `Developing the Provider` commands in `bash` because `sh` scrips are invoked as part of these.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.15+ is **required**). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

First clone the repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-azurerm`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-azurerm
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-azurerm
```

Once inside the provider directory, you can run `make tools` to install the dependent tooling required to compile the provider.

At this point you can compile the provider by running `make build`, which will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-azurerm
...
```

You can also cross-compile if necessary:

```sh
GOOS=windows GOARCH=amd64 make build
```

In order to run the `Unit Tests` for the provider, you can run:

```sh
$ make test
```

The majority of tests in the provider are `Acceptance Tests` - which provisions real resources in Azure. It's possible to run the entire acceptance test suite by running `make testacc` - however it's likely you'll want to run a subset, which you can do using a prefix, by running:

```sh
make acctests SERVICE='resource' TESTARGS='-run=TestAccAzureRMResourceGroup' TESTTIMEOUT='60m'
```

The following Environment Variables must be set in your shell prior to running acceptance tests:

- `ARM_CLIENT_ID`
- `ARM_CLIENT_SECRET`
- `ARM_SUBSCRIPTION_ID`
- `ARM_TENANT_ID`
- `ARM_ENVIRONMENT`
- `ARM_METADATA_HOST`
- `ARM_TEST_LOCATION`
- `ARM_TEST_LOCATION_ALT`
- `ARM_TEST_LOCATION_ALT2`

**Note:** Acceptance tests create real resources in Azure which often cost money to run.

---

## Developer: Instructing Terraform to use your locally compiled Azure Provider binary

After successfully compiling the Azure Provider, you must [instruct Terraform to use your locally compiled provider binary]https://www.terraform.io/docs/commands/cli-config.html#development-overrides-for-provider-developers instead of the official binary from the Terraform Registy.

For example, add the following to `~/.terraformrc` for a provider binary located in `/home/developer/go/bin`:

```hcl
provider_installation {

  # Use /home/developer/go/bin as an overridden package directory
  # for the hashicorp/azurerm provider. This disables the version and checksum
  # verifications for this provider and forces Terraform to look for the
  # azurerm provider plugin in the given directory.
  dev_overrides {
    "hashicorp/azurerm" = "/home/developer/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

---

## Developer: Generating Resource ID Formatters, Parsers and Validators

You can generate a Resource ID Formatter, Parser and Validator by adding the following line to a `resourceids.go` within each Service Package (for example `./azurerm/internal/services/someservice/resourceids.go`):

```go
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Server -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AnalysisServices/servers/Server1
```

Where `name` is the name of the Resource ID Type - and `id` is an example Resource ID with placeholder data.

When `make generate` is run, this will then generate the following for this Resource ID:

* Resource ID Struct, containing the fields and a Formatter to convert this into a string - and the associated Unit Tests.
* Resource ID Parser (`./parse/{name}.go`) - to be able to parse a Resource ID into said struct - and the associated Unit Tests.
* Resource ID Validator (`./validate/{name}_id.go`) - to validate the Resource ID is what's expected (and not for a different resource) - and the associated Unit Tests.

---

## Developer: Scaffolding the Website Documentation

You can scaffold the documentation for a Data Source by running:

```sh
$ make scaffold-website BRAND_NAME="Resource Group" RESOURCE_NAME="azurerm_resource_group" RESOURCE_TYPE="data"
```

You can scaffold the documentation for a Resource by running:

```sh
$ make scaffold-website BRAND_NAME="Resource Group" RESOURCE_NAME="azurerm_resource_group" RESOURCE_TYPE="resource" RESOURCE_ID="/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1"
```
