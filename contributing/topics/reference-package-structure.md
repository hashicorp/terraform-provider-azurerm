# Package and Folder Structure

Each Azure service has its own package under `internal/services/<service>/`. Within a service, every resource, data source, list resource and action lives in its own sub-package.

## Layout

```text
internal/services/<service>/
├── client/                  # SDK clients for the service
├── parse/                   # Resource ID formatters and parsers
├── validate/                # Validation functions, including Resource ID validators
├── registration.go          # Service Registration
├── resourceids.go           # Used to generate Resource ID formatters, parsers and validators
├── resources/
│   └── <resource_name>/     # e.g. subnet
│       ├── resource.go
│       ├── r_schema.go
│       ├── r_create.go
│       ├── r_read.go
│       ├── r_update.go
│       ├── r_delete.go
│       ├── data_source.go
│       ├── list.go
│       ├── migration/       # State migrations for this resource
│       ├── resource_test.go
│       ├── data_source_test.go
│       └── list_test.go
└── actions/
    └── <action_name>/       # e.g. flush_databases
        ├── action.go
        └── action_test.go
```

Only the files a resource needs are present - for example, a resource without a data source has no `data_source.go`.

## Naming

* The folder is the Terraform name without the `azurerm_` prefix - `azurerm_subnet` is `resources/subnet/`, `azurerm_network_security_group` is `resources/network_security_group/`.
* Actions follow the same pattern under `actions/`.
* The package uses the folder name (`package subnet`), and the tests use the same name with a `_test` suffix (`package subnet_test`).

See [Naming](reference-naming.md) for how to name the resource itself.

## What Goes in Each File

| File | Contents |
| ---- | -------- |
| `resource.go` | The Resource struct and its entry points (e.g. `ResourceType()`). |
| `r_schema.go` | The schema: arguments, attributes and their validation. |
| `r_create.go` | Create. |
| `r_read.go` | Read. |
| `r_update.go` | Update. |
| `r_delete.go` | Delete. |
| `r_expand.go` / `r_flatten.go` | Converting between the Terraform configuration/state and the API payloads, where the resource needs them. |
| `data_source.go` | The Data Source for this resource, if there is one. |
| `list.go` | The List Resource (see [Adding a new List Resource](guide-list-resource.md)). |
| `resource_test.go` | Acceptance tests for the Resource. |
| `data_source_test.go` | Acceptance tests for the Data Source. |
| `list_test.go` | Acceptance tests for the List Resource. |
| `resource_identity_gen_test.go` | Resource Identity tests (see [Adding Resource Identity](guide-resource-identity.md)). |

## Shared Code in the Service Package

Anything used by more than one resource stays at the root of the Service Package:

* `client/` - the Client for the service.
* `parse/` and `validate/` - Resource ID parsing and validation.
* `registration.go` - every resource, data source and action in the service is registered here.
* `resourceids.go` - the Resource ID definitions used for generation.

> **Note:** `resourceids.go` is a legacy requirement which will be removed in the future. For now it is still needed in services that use generated Resource IDs.

## State Migrations

State migrations for a resource live in a `migration` folder inside that resource's sub-package (e.g. `resources/app_configuration_key/migration/`). See [Adding State Migrations](guide-state-migrations.md).

## Adding a New Resource

1. Create `resources/<resource_name>/` in the Service Package.
2. Add the files from [What Goes in Each File](#what-goes-in-each-file) that your resource needs.
3. Register the resource in the Service Package's `registration.go`.
4. Add the acceptance tests in the same folder, using the `<resource_name>_test` package.

See [Adding a new Resource](guide-new-resource.md) for the full steps.

## Resources That Can Only Exist Once

Some Azure resources can only exist once per region or subscription, so their tests can't run at the same time as each other. The tests for these live together with the top-level resource - see [Sequential Acceptance Tests](reference-acceptance-testing.md#sequential-acceptance-tests).

## Running the Tests

Run all of the tests for a service with:

```sh
make acctests SERVICE='<service>' TESTTIMEOUT='60m'
```

Or a subset by name with `TESTARGS`:

```sh
make acctests SERVICE='<service>' TESTARGS='-run=TestAccSubnet_basic' TESTTIMEOUT='60m'
```

See [Running the Tests](running-the-tests.md) for the environment variables needed.
