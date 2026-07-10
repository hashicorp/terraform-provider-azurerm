# tfpdk

Terraform Provider Development Kit

DISCLAIMER: The author primarily works on the [AzureRM Provider](https://github.com/hashicorp/terraform-provider-azurerm) so this tool is *heavily* biased towards the code styles and requirements of that provider currently.  Future development will look to genericise the output (eventually 😜)

NOTE: Expects to be run from the root of a validly named Terraform provider e.g. `./terraform-provider-myprovider`

NOTE: Relies on a locally installed terraform binary and appropriately configured [dev overrides file](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers). This
allows the Terraform binary to skip the `init` step by informing it where your compiled provider binary can be found. (you still currently need a local .tf file with your [provider configured](https://www.terraform.io/docs/language/providers/configuration.html) though, sorry!)

Run from the root of the provider project.

## TODO

- [x] untyped resource template
- [x] `go fmt` outputs
- [x] Add new resources / data sources to appropriate registration
- [x] typed and untyped Data Sources
- [x] init a new provider - git clone [scaffold](https://github.com/hashicorp/terraform-provider-scaffolding)?
- [ ] optionally populate the Typed SDK into the init step?
- [x] Dynamically read provider name for item generation
- [ ] Populate `IDValidationFunc()` in template for IDs (Pandora)
- [ ] Clients?
- [ ] Autocomplete CLI?
- [x] Spike doc generation
- [ ] Extend doc gen to deal with Blocks
- [ ] migrate to ~~plugin sdk~~ something else to read schema directly for doc gen as Terraform's JSON output lacks necessary detail. (e.g. timeouts, ForceNew flags etc)
- [x] update commands and templates to allow non-hashicorp providers ~~and sources other than github~~ (e.g. import paths)

## Commands

Usage: `tfpdk [--version] [--help] <command> [<args>]`

```
Available commands are:
config            Generate a local config file for common options.
document          generates documentation from a resource.
generate          generates a typed resource (and optionally a list resource and/or data source) from the Pandora Data API.
init              initialises a new provider based on the scaffold project.
servicepackage    Creates a directory for a new Service Package and scaffolds out the basics to use it.
```

## Usage Examples

### Generate a typed resource from the Pandora Data API

```shell
tfpdk generate -arm-type "Microsoft.RedHatOpenShift/openShiftClusters" -name redhat_openshift_cluster -go-name RedHatOpenShiftCluster -servicepackage redhatopenshift
```

Will create `{providername}/internal/services/redhatopenshift/redhat_openshift_cluster_resource.go` with the schema, model structs, CRUD and expand/flatten functions populated from the API definitions (read from a running Pandora Data API).

Add `-list` to also generate a list resource, `-data-source` to also generate a data source, or `-gen-resource=false -list` to generate only the list resource for a resource that already exists.

### Document your new resource (relies on the `Description` fields in your schema, so populate those for best effect)

```shell
tfpdk document -type resource -name ShinyNewService -id "00000000-0000-0000-0000-000000000000"
```

### Generate a Local config to set options and paths

```shell
tfpdk config
```

will create `.tfpdk.hcl` in the root of the provider, which can be updated to configure provider specifics that may differ from the defaults. (by "defaults", I mean what we use on AzureRM), example output:

```
provider_name                   = "azurerm"
service_packages_path           = "internal/services"
provider_github_org             = "hashicorp"
docs_path                       = "docs"
resource_docs_directory_name    = "r"
data_source_docs_directory_name = "d"
use_typed_sdk                   = false
```

**Note:** Command line options override values from the config file.

## Command Documentation

`tfpdk [command]`
- [init](docs/init.md)
- [config](docs/config.md)
- [servicepackage](docs/servicepackage.md)
- [document](docs/document.md)
