# Terraform Provider for Azure (Resource Manager)

The AzureRM Provider supports Terraform 0.10.x and later - but Terraform 0.12.x is recommended.

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
  version = "=1.38.0"

  # More information on the authentication methods supported by
  # the AzureRM Provider can be found here:
  # http://terraform.io/docs/providers/azurerm/index.html

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
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["10.0.0.0/16"]
}
```

Further [usage documentation is available on the Terraform website](https://www.terraform.io/docs/providers/azurerm/index.html).

## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.10.x +
* [Go](https://golang.org/doc/install) version 1.13.x (to build the provider plugin)

If you're on Windows you'll also need:
* [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)
* [Git Bash for Windows](https://git-scm.com/download/win)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.*

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".*

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is **required**). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

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

In order to run the Unit Tests for the provider, you can run:

```sh
$ make test
```

The majority of tests in the provider are Acceptance Tests - which provisions real resources in Azure. It's possible to run the entire acceptance test suite by running `make testacc` - however it's likely you'll want to run a subset, which you can do using a prefix, by running:

```
make testacc TESTARGS='-run=TestAccAzureRMResourceGroup'
```

The following Environment Variables must be set in your shell prior to running acceptance tests:

- `ARM_CLIENT_ID`
- `ARM_CLIENT_SECRET`
- `ARM_SUBSCRIPTION_ID`
- `ARM_TENANT_ID`
- `ARM_ENVIRONMENT`
- `ARM_TEST_LOCATION`
- `ARM_TEST_LOCATION_ALT`
- `ARM_TEST_LOCATION_ALT2`

**Note:** Acceptance tests create real resources in Azure which often cost money to run.

