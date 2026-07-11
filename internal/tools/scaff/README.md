# Scaff

## Description

`scaff` is a code-scaffolding CLI for the provider. It reads resource schemas
from a running [Pandora Data API](#running-the-pandora-data-api) and/or the
existing provider source and generates idiomatic Go, including:

- **New resources & data sources** — schema, CRUD, and expand/flatten functions
  for the Typed SDK (`internal/sdk`).
- **List resources** — makes an existing resource list-ready (Resource Identity
  - a reusable flatten method) and generates the `*_resource_list.go` (and its
  acceptance test), registering it in the service package.
- **Service package scaffolding**, and **documentation**.

It exposes the following commands:

| Command          | Purpose                                                                             |
| ---------------- | ----------------------------------------------------------------------------------- |
| `generate`       | Generate a typed resource (and optionally a data source / list resource).           |
| `upgrade`        | Upgrade an existing resource to support List (adds identity + flatten).             |
| `servicepackage` | Create a new service package directory and scaffold `registration.go` + `client/`.  |
| `document`       | Generate provider documentation for a resource or data source.                      |
| `config`         | Write a local `.scaff.hcl` config file with defaults for the options below.         |

## Set-up

Run every command from the **provider repository root** (`terraform-provider-azurerm`).

Invoke it with `go run` (no build step required):

```sh
go run ./internal/tools/scaff <command> [options]
```

or build a binary once:

```sh
go build -o bin/scaff ./internal/tools/scaff
./bin/scaff <command> [options]
```

Pass `-h`/`-help` to any command to see its full option list.

Most commands read schema from a Pandora Data API running at
`http://localhost:8080` — see [Running the Pandora Data API](#running-the-pandora-data-api).
Defaults can be overridden per-invocation with flags, or persisted in a
[`.scaff.hcl`](#configuration-scaffhcl) file.

## Usage

### Generate

Generates a typed (`internal/sdk`) resource — schema, CRUD, and expand/flatten
functions — from the API definitions served by the Pandora Data API. Existing
files are left untouched unless `-overwrite` is passed.

#### Prerequisites

A running Pandora Data API reachable at `schema_api_url` (default
`http://localhost:8080`). See [Running the Pandora Data API](#running-the-pandora-data-api).

#### Options

| Flag              | Required | Description                                                                         |
| ----------------- | :------: | ----------------------------------------------------------------------------------- |
| `-arm-type`       |   yes    | ARM resource type, e.g. `Microsoft.RedHatOpenShift/openShiftClusters`.              |
| `-name`           |   yes    | Terraform resource name, snake_case without the provider prefix.                    |
| `-go-name`        |          | Go identifier for the resource; derived from `-name` when omitted.                  |
| `-servicepackage` |          | Service package directory; derived from the service name when omitted.              |
| `-api-version`    |          | API version; defaults to the latest non-preview version.                            |
| `-service`        |          | Explicit Pandora service name; overrides the value derived from `-arm-type`.        |
| `-resource`       |          | Explicit Pandora resource key; overrides the value derived from `-arm-type`.        |
| `-list`           |          | Also generate a list resource and its acceptance test.                             |
| `-gen-resource`   |          | Generate the resource file (default `true`); set `false` with `-list` for list-only.|
| `-data-source`    |          | Also generate a data source (shares nested block structs with the resource).        |
| `-path`           |          | Output directory; defaults to `{service_packages_path}/{servicepackage}`.           |
| `-pandora-url`    |          | Pandora Data API base URL; defaults to the configured `schema_api_url`.             |
| `-provider`       |          | Provider name override, e.g. `azurerm`.                                             |
| `-org`            |          | Provider GitHub org override, e.g. `hashicorp`.                                     |
| `-input`          |          | HCL file describing one or more resources; per-resource flags are ignored when set. |
| `-overwrite`      |          | Overwrite existing generated files.                                                 |

#### Examples

```sh
# A single resource:
go run ./internal/tools/scaff generate \
  -arm-type="Microsoft.RedHatOpenShift/openShiftClusters" \
  -name="redhat_openshift_cluster" \
  -go-name="RedHatOpenShiftCluster" \
  -servicepackage="redhatopenshift"

# Resource + data source:
go run ./internal/tools/scaff generate ... -data-source

# List resource only (the resource already exists with a flatten method):
go run ./internal/tools/scaff generate ... -gen-resource=false -list

# A batch of resources described in HCL:
go run ./internal/tools/scaff generate -input="internal/tools/scaff/examples/generate.hcl"
```

See [`examples/generate.hcl`](examples/generate.hcl) for the input-file format.

### Upgrade

Upgrades an existing resource (typed or untyped/native Plugin SDK) so it can
support List. It adds Resource Identity and refactors `Read` into a reusable
flatten method when those are missing, then — with `-list` (the default) —
generates the list resource, its acceptance test, and registers it in the
service package's `ListResources()`.

Most metadata (SDK package, read model, and list operations) is derived from the
resource source and the vendored `go-azure-sdk`, so a running Data API is only
needed as a fallback when the list operation cannot be resolved from source.
Supply `-arm-type` (or `-service` + `-resource`), or `-read-model`, in that case.

The command performs a **dry run by default**, printing the proposed changes;
pass `-write` to apply them.

#### Options

| Flag                    | Required | Description                                                                        |
| ----------------------- | :------: | ---------------------------------------------------------------------------------- |
| `-file`                 |   yes    | Path to the existing resource `.go` file to upgrade.                               |
| `-list`                 |          | Make the resource list-ready and generate the list resource (default `true`).      |
| `-identity`             |          | Add Resource Identity if missing.                                                  |
| `-flatten`              |          | Refactor `Read` into a reusable flatten method if missing.                         |
| `-write`                |          | Write changes to disk; without it the command performs a dry run.                  |
| `-overwrite`            |          | Overwrite an existing generated list file.                                         |
| `-arm-type`             |          | ARM type; used to resolve list operations / read model when not source-derivable.  |
| `-service`, `-resource` |          | Explicit Pandora service name / resource key.                                       |
| `-api-version`          |          | API version; defaults to the latest non-preview version.                            |
| `-read-model`           |          | SDK read-model type name; overrides the value resolved from Pandora.               |
| `-list-method`          |          | Parent-scoped SDK list method name (without the `Complete` suffix).                 |
| `-name`, `-go-name`     |          | Terraform resource name / Go identifier; derived from the file when omitted.        |
| `-identity-properties`  |          | `-properties` value for the identity test `go:generate` (default `name,resource_group_name`). |
| `-pandora-url`          |          | Pandora Data API base URL; defaults to the configured `schema_api_url`.             |
| `-provider`, `-org`     |          | Provider name / GitHub org overrides.                                               |
| `-input`                |          | HCL file describing one or more resources; per-resource flags are ignored when set. |

#### Examples

```sh
# Report what an upgrade would do (read-only dry run):
go run ./internal/tools/scaff upgrade \
  -file internal/services/monitor/monitor_workspace_resource.go

# Make a resource list-ready and write the changes:
go run ./internal/tools/scaff upgrade \
  -file internal/services/monitor/monitor_workspace_resource.go -write

# Upgrade a batch of resources described in HCL:
go run ./internal/tools/scaff upgrade -input="internal/tools/scaff/examples/upgrade-network.hcl"
```

In an HCL input file `list` defaults to `true`, so a block need only set
`list = false` to opt out. See [`examples/upgrade.hcl`](examples/upgrade.hcl) and
[`examples/upgrade-network.hcl`](examples/upgrade-network.hcl).

> **Note:** when running against real service directories, keep `write = false`
> (dry run) until you've reviewed the output. Writing applies changes in place.

### Service Package

Creates a new service package directory under `service_packages_path` and
scaffolds the basics needed to use it: `registration.go` and `client/client.go`.

| Flag              | Required | Description                                     |
| ----------------- | :------: | ----------------------------------------------- |
| `-servicepackage` |   yes    | The name of the service package.                |
| `-typed`          |          | Use the Typed SDK (defaults to `use_typed_sdk`).|

```sh
go run ./internal/tools/scaff servicepackage -servicepackage="redhatopenshift"
```

### Document

Generates provider documentation (markdown) for a resource or data source from
the schema API.

| Flag              | Description                                                        |
| ----------------- | ----------------------------------------------------------------- |
| `-name`           | The name of the resource.                                         |
| `-servicepackage` | The service package the resource or data source belongs to.       |
| `-type`           | The item to document: `resource` or `data-source`.                |
| `-id`             | Example resource ID (only valid for `-type=resource`).            |
| `-schemaapiurl`   | The schema API URL (default `http://localhost:8080`).             |

```sh
go run ./internal/tools/scaff document \
  -name="redhat_openshift_cluster" -servicepackage="redhatopenshift" -type="resource"
```

### Config

Writes a `.scaff.hcl` file in the current directory pre-filled with the current
defaults, which you can then edit.

```sh
go run ./internal/tools/scaff config
```

## Configuration (`.scaff.hcl`)

When a `.scaff.hcl` file is present in the working directory, its values
override the built-in defaults for the options below. All keys are optional.

| Key                               | Default                | Description                                              |
| --------------------------------- | ---------------------- | -------------------------------------------------------- |
| `provider_name`                   | *(from `go.mod`)*      | Provider name, e.g. `azurerm`.                           |
| `provider_canonical_name`         |                        | Canonical provider name.                                 |
| `provider_github_org`             | `hashicorp`            | Provider GitHub organisation.                            |
| `service_packages_path`           | `internal/services`    | Directory containing the service packages.               |
| `schema_api_url`                  | `http://localhost:8080`| Base URL of the Pandora Data API.                        |
| `docs_path`                       | `docs`                 | Directory documentation is written to.                   |
| `docs_version`                    | `legacy`               | Docs layout: `legacy` or `registry`.                     |
| `resource_docs_directory_name`    | `r`                    | Sub-directory for resource docs.                         |
| `data_source_docs_directory_name` | `d`                    | Sub-directory for data source docs.                      |
| `use_typed_sdk`                   | `true`                 | Whether generated code targets the Typed SDK.            |

## Running the Pandora Data API

The `generate` and `document` commands (and `upgrade` as a fallback) resolve
schemas from a Pandora Data API. The helper script clones (or reuses) the
Pandora repository, builds the Data API, and serves it on `http://localhost:8080`:

```sh
internal/tools/scaff/scripts/serve-data-api.sh
```

Useful options (see `--help` for the full list):

```sh
internal/tools/scaff/scripts/serve-data-api.sh --port 8099          # custom port
internal/tools/scaff/scripts/serve-data-api.sh --services Network   # only load some services
internal/tools/scaff/scripts/serve-data-api.sh --update             # git pull an existing checkout first
```

It reuses a sibling `pandora/` checkout when present, otherwise clones into
`~/.cache/scaff/pandora`. Override the location and source with the `PANDORA_DIR`
and `PANDORA_REPO_URL` environment variables.
