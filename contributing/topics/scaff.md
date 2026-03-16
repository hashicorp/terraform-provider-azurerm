<!-- Copyright IBM Corp. 2014, 2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

# Provider Scaffolding (scaff)

`scaff` is a Terraform AzureRM Provider scaffolding command line tool.
It generates resource, data source, list resource, and documentation files along with test files which adhere to the latest best practices.
These files are heavily commented with instructions, serving as the best way to get started with provider development.

## Overview workflow steps

1. Install `scaff`.
1. Use `scaff` to generate provider code.
1. Go through the generated code, customizing as necessary.
1. Run, test, refine.
1. Remove "TIP" comments.
1. Submit a pull request.

## Running `scaff`

1. Clone the [Terraform AzureRM Provider](https://github.com/hashicorp/terraform-provider-azurerm) repository.
1. Install `scaff`.

    ```sh
    make scaff
    ```

1. Change into the appropriate directory.
    * For resources and list resources, this is the service directory where the new entity will reside, e.g. `internal/services/privatedns`.
    * For functions, this is `internal/functions`.
1. Generate the resource, list resource, data source or function. For example,
    * `scaff resource -name="fw_mssql_database" -service_package_name="MSSQL" ...`.
    * `scaff list-resource -service="PrivateDns" -resource="CNameRecord" ...`.
    * `scaff list-documentation -path=internal/services/network`.

To get help, enter `scaff` without arguments.

## Usage

### Help

```console
scaff --help
```

```
Usage:
  scaff [command]

Available Commands:
  list-documentation Generate documentation for list resources
  list-resource      Generate list resource boilerplate code
  resource           Create scaffolding for a resource

Flags:
  -h, --help   help for scaff
```

### List Resource

Generate list resource boilerplate code from JSON file or CLI arguments.

```console
scaff list-resource --help
```

```
Usage: scaff list-resource [options]

  Generates list resource boilerplate code. Can accept either a JSON file
  with multiple resource definitions or individual CLI arguments for a single resource.

Options:
  -json=<path>                    Path to JSON file containing resource definitions (use instead of full cli)
  -service=<name>                 (Required) Service name (e.g., "PrivateDns")
  -resource=<name>                (Required) Resource name (e.g., "CNameRecord")
  -include_service_in_name=<bool> (Optional) Include service name in generated identifiers, defaults to false
  -full_parent=<name>             (Optional) Full parent resource name
  -parent=<name>                  (Optional) Parent resource name (e.g., "PrivateDnsZone"), defaults to resource_group
  -parent_terraform_name=<name>   (Optional) Parent Terraform resource name (e.g., "private_dns_zone")
  -terraform_name=<name>          (Optional) Terraform resource name (e.g., "private_dns_cname_record")
  -id_structure=<type>            (Required) ID structure type (e.g., "privatedns.RecordType")
  -path=<path>                    (Optional) Output path for generated files


Examples:
  # Using JSON file
  scaff list-resource -json=internal/tools/scaff/commands/input_example/listResources.json

  # Using CLI arguments
  scaff list-resource -service="PrivateDns" -resource="CNameRecord" -parent="PrivateDnsZone" \
    -terraform_name="private_dns_cname_record" -id_structure="privatedns.RecordType" \
    -path="internal/services/privatedns/"
```

### List Documentation

Generate documentation for list resources by scanning Go source files.
> [!NOTE]
      This command can only be used when the *_resource_list.go files already exist.

```console
scaff list-documentation --help
```

```
Usage: scaff list-documentation [options]

  Generates documentation for list resources by scanning Go source files
  for files ending with '_resource_list.go'.

Options:
  -path=<path>              (Required) Path to file or directory to scan (required)
  -subcategory=<name>       (Optional) Override the subcategory/section (e.g., "Network", "Database")
  -addsectiontoname         (Optional) Prepend section name to attribute names (boolean flag)

Examples:
  # Basic usage
  scaff list-documentation -path=internal/services/network

  # With subcategory override
  scaff list-documentation -path=internal/services/mssql -subcategory=Database

  # With section name prefixing enabled
  scaff list-documentation -path=internal/services/mssql -addsectiontoname
```

### Resource

Create scaffolding for a resource using the Terraform Plugin Framework.

```console
scaff resource --help
```

```
Usage: scaff resource -name "some_resource_name" -service_package_name="someservice" \
  -rp_name="sql" -client_name="SomeClient" [-updatable=true] [-no_resource_group=true] \
  -api_version="2023-08-01-preview" -id_type="commonids.SqlDatabaseId" [-sdk_name="databases"] \
  -id_segments="SubscriptionId,ResourceGroupName,ServerName,DatabaseName" \
  [-uses_lro_crud=true] [-use_create_options=true] [-use_read_options=true] \
  [-use_update_options=true] [-use_delete_options=true]

Parameters:
  -name (Required)                 The name of the resource to scaffold
  -service_package_name (Required) The name of the service package to scaffold the resource into
  -rp_name (Required)              The name of the resource provider of the new resource
  -client_name (Required)          The name of the client used to manage the new resource
  -api_version (Required)          The API version of the resource to scaffold (e.g., 2025-01-01)
  -id_type (Required)              The type of resource to scaffold (e.g., 'commonids.AppServiceId',
                                   or 'virtualmachines.VirtualMachineId')
  -id_segments (Required)          The User-Specified Segment names for the ID, order matters

  -updatable (Optional)            Whether the new resource can be updated (i.e., any schema property
                                   is not going to be 'ForceNew')
  -uses_lro_crud (Optional)        The new resource uses LROs for Create, Update, and Delete
  -use_create_options (Optional)   The new resource uses OperationOptions for Create
  -use_read_options (Optional)     The new resource uses OperationOptions for Read
  -use_update_options (Optional)   The new resource uses OperationOptions for Update
  -use_delete_options (Optional)   The new resource uses OperationOptions for Delete
  -config_validators (Optional)    Does the resource have configuration validators
  -no_resource_group (Optional)    Set to true if the resource is not created in a resource group,
                                   or if the RG is inferred from a parent resource ID
  -sdk_name (Optional)             The name of the SDK used to manage the new resource. If omitted,
                                   the first slug of the id_type value will be used

Example:
  scaff resource -name="fw_mssql_database" -service_package_name="MSSQL" -rp_name="sql" \
    -client_name="DatabasesClient" -updatable=true -no_resource_group=true \
    -api_version="2023-08-01-preview" -id_type="commonids.SqlDatabaseId" -sdk_name="databases" \
    -id_segments="SubscriptionId,ResourceGroupName,ServerName,DatabaseName" \
    -uses_lro_crud=true -use_read_options=true
```

### JSON Configuration

For batch generation of multiple list resources, you can use a JSON configuration file:

```json
[
  {
    "service": "PrivateDns",
    "resource": "AaaaRecord",
    "include_service_in_name": true,
    "full_parent": "",
    "parent": "PrivateDnsZone",
    "terraform_name": "private_dns_aaaa_record",
    "id_structure": "privatedns.RecordType",
    "path": "internal/services/privatedns/"
  }
]
```

Then run:

```sh
scaff list-resource -json=path/to/config.json
```
