# Glossary

This document contains a summary of the terminology used within the Azure Provider.

### Azure Resource ID

An Azure Resource ID is used to uniquely identify this Resource within Azure - in almost all cases this is a Path of Key-Value Pairs, for example:

> /subscriptions/11112222-3333-4444-555566667777/resourceGroups/myGroup

Contains the Key-Value pairs:

> `subscriptions`: `11112222-3333-4444-555566667777`
> `resourceGroups`: `myGroup`

As the Azure Resource ID is comprised of user-specified Key-Value Pairs, the Azure Resource ID is predictable.

### Data Plane API

A Data Plane API provides access to data for resources provisioned via the Resource Manager API. Some examples:

* The App Configuration Data Plane API allows for managing Keys and Features within an App Configuration.
* The Storage Data Plane API allows for the uploading/downloading of Blobs within a Storage Container (within a Storage Account).

### Embedded SDK

An Embedded SDK is an SDK that has been added directly into the providers code base (usually into `services/{name}/sdk`) rather than using go modules and vendoring it into `/vendor`.

Whilst we generally vendor SDKs instead, we have a number of SDKs which aren't available elsewhere and are instead vendored into the codebase (see [High Level Overview](high-level-overview.md) for more information).

### Resource ID Formatter

A Resource ID Formatter is a Resource ID Struct which implements the `ID()` method - returning the (Azure) Resource ID as a string - which must be parseable using the associated Resource ID Parser.

These are generally (but not always) auto-generated - see Terraform Managed Resource ID’s below for more information.

### Resource ID Parser

A Resource ID Parser parses an (Azure) Resource ID into a Resource ID Struct - generally case-sensitively (since both Terraform Core and some downstream Azure API’s are case sensitive), but optionally case-insensitively where required.

These are generally (but not always) auto-generated - see Terraform Managed Resource ID’s below for more information.

### Resource ID Struct

A Resource ID Struct is a Golang Struct defining the user-specifiable values within an (Azure) Resource ID. For example, in the case of a Resource Group ID that would be the Subscription ID and Resource Group name.

A Resource ID Struct should have an associated Resource ID Formatter, Parser and (optionally) Validator.

These are generally (but not always) auto-generated - see Terraform Managed Resource ID’s below for more information.

### Resource ID Validator

A Resource ID Validator is a Terraform Validation function which validates that the specified value is a Resource ID of the expected Type (for example a Subnet ID validator checks it’s a Subnet ID).

The value is parsed case-sensitively (in some cases, an optional case-insensitive validation function is also available) using the associated Resource ID Parser.

This Resource ID Validator can then be used as a validation function within Terraform Schema fields as necessary - to confirm that the user-specified value (for example, for a Subnet ID) is actually the specified type (for example, a Subnet ID) and not another Resource ID or value (for example, a Virtual Network ID).

These are generally (but not always) auto-generated - see Terraform Managed Resource ID’s below for more information.

### Resource Manager API

> Some Service Teams refer to this as "Management Plane".

A Resource Manager API is used to provision resources within an Azure Subscription/Management Group, for example a Resource Group or a Virtual Machine.

Whilst the Resource Manager API can be used to provision resources, resources within those are generally exposed via Data Plane APIs (see above) - for example Blobs within a Storage Account.

### Service Package

A Service Package is a grouping of Data Sources and Resources (and any other associated functionality) which are related together, for example `Cosmos` or `Compute`.

Each Service Package contains a Service Registration which defines the Data Sources and Resources available within that Service Package.

Whilst these tend to map 1:1 to Azure Resource Providers (for example the `cosmos` Service Package contains the CosmosDB resources) - some are intentionally split out where the Resource Provider (or Service Package) would otherwise be too large (for example the Network package has Load Balancers split out).

### Service Registration

Each Service Package contains a Service Registration which defines the Data Sources and Resources available within that Service Package.

This is either a Typed Service Registration or an Untyped Service Registration (documented below) - both available within the Typed Plugin SDK.

Note that a Service Registration can be both a Typed and Untyped Service Registration by implementing both the Typed and Untyped Service Registration interfaces. This allows the mixing of both Typed and Untyped Data Sources and Resources within a Service Package.

### State Migration

A State Migration is used when a resource has been changed to expect something different in the state than what previous version of the provider have written to it. An example of this is if Azure started to return a Resource ID value in a different case. rather than showing this during the plan, we can write a state migration to update the ID values transparently with no action required by a user. These are found in `services/service/migrations` and documentation on how to write them can be found in the [Terraform Plugin SDK](https://www.terraform.io/plugin/sdkv2/resources/state-migration) documentation.

### Terraform Managed Resource ID

A Terraform Managed Resource ID is a Resource ID defined in Terraform, rather than set by the Remote API.

The Azure Provider is moving to use Terraform Managed Resource ID’s for all resources, since these are known ahead of time - which avoids issues with API’s changing these Resource ID’s over time (either in casing, or renaming segments altogether).

At present these are defined in a `resourceids.go` file within each Service Package, which generates a Resource ID Formatter, Parser and Validator for this Resource ID.

### Terraform Resource Data

Terraform Resource Data is a wrapper around the values within either the Terraform Configuration/State, depending on when this is called.

Values within the Resource Data can be accessed using `d.Get` (for example `d.Get(“some_field”).(string)`) and set using `d.Set` (for example `d.Set(“some_field”, “hello”)`.

### Terraform Resource ID

Each Data Source and Resource within Terraform has a Resource ID used to keep track of this resource, set at creation/import time.

For a Resource this is set in the Create function after the resource has been successfully provisioned (or at Import time, when imported) - and then used in the Delete, Read and Update functions to look up this resource.

Since Data Sources look up information about existing resources - and as such don’t have a Create method - these instead set the Resource ID within the Read function.

### Typed Data Source

A Typed Data Source is a Terraform Data Source built using the Typed Plugin SDK, allowing this Data Source to be defined using Native Go Types.

### Typed Resource

A Typed Resource is a Terraform Resource built using the Typed Plugin SDK, allowing this Resource to be defined using Native Go Types.

### Typed Plugin SDK

The Typed Plugin SDK is an abstraction over the Terraform Plugin SDK housed within the AzureRM Provider repository - which allows Terraform Data Sources and Resources to be built using Native Go Types.

The Typed Plugin SDK contains both Golang Interfaces for Data Sources and Resources (which allows verifying these are valid at compile-time) - and a wrapper around Terraform Resource Data which allows for values from the Terraform Configuration to be  Serialized/Deserialized into a Native Go Struct.

More information can be found in [the documentation for the Typed Plugin SDK](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/internal/sdk).

### Typed Service Registration

A Typed Service Registration returns a list of the Typed Data Sources and Typed Resources which are available within that Service Package.

This is implemented within [the Typed Plugin SDK](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/internal/sdk) as the interface `TypedServiceRegistration` (see also: `TypedServiceRegistrationWithAGitHubLabel`).

### Untyped Data Source

An Untyped Data Source is a Terraform Data Source built using the Terraform Plugin SDK directly, which looks up information about an existing Resource. These are exposed as a function which returns an instance of the Plugin SDK’s `Resource` struct - implementing whichever methods are necessary (generally, the Schema and Read/Timeouts functions).

The Terraform Resource Data can be used to set fields into the Terraform State - and to set the ID using `d.SetId(“”)`.

### Untyped Resource

An Untyped Resource is a Terraform Resource built using the Terraform Plugin SDK directly, which manages this Resource (through either creation/import onwards). These are exposed as a function which returns an instance of the Plugin SDK’s `Resource` struct - implementing whichever methods are necessary (generally, the Schema and Create/Read/Update/Delete/Import/Timeouts functions).

The Terraform Resource Data can be used to retrieve fields from the Terraform Configuration/set fields into the Terraform State - and to get/set the ID using `d.Id()` / `d.SetId(“”)`.

### Untyped Service Registration

An Untyped Service Registration returns a list of the Untyped Data Sources and Untyped Resources which are available within that Service Package.

This is implemented within [the Typed Plugin SDK](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/internal/sdk) as the interface `UntypedServiceRegistration` (see also: `UntypedServiceRegistrationWithAGitHubLabel`).
