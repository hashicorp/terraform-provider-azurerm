## Tool: `update-go-azure-sdk`

This tool updates the version of `hashicorp/go-azure-sdk` which is vendored in this repository, and then subsequently tries to update the Provider to account for any new API Versions.

Whilst this isn't going to be infallible, the intention is that when a new API Version is added to `hashicorp/go-azure-sdk` (via `hashicorp/pandora`) then the AzureRM Provider can be automatically updated to use that API Version - with any errors output for human review.

This tool makes a number of assumptions about the repository design, namely that an `internal` folder exists and that `hashicorp/go-azure-sdk` is not referenced outside of this directory - and in addition also tool assumes that Git is available and pre-configured.

Whilst this tool is intended to primarily be run in automation - [the `update-api-version` tool](../update-api-version) is better suited for interactive usage (since it won't perform rollbacks, and assumes the `hashicorp/go-azure-sdk` dependency is up-to-date).

### Example Usage

```bash
 $ go build . && ./update-go-azure-sdk --new-sdk-version=v0.20231005.1153009 --azurerm-repo-path=../../../                                                                                                      (tool/update-go-azure-sdk ⚡)
2023-10-11T10:47:53.665+0200 [INFO]  New SDK Version is "v0.20231005.1153009"
2023-10-11T10:47:53.665+0200 [INFO]  The `hashicorp/terraform-provider-azurerm` repository is located at "../../../"
2023-10-11T10:47:53.665+0200 [INFO]  A path to the `hashicorp/go-azure-sdk` repository was not provided - will clone on-demand
2023-10-11T10:47:53.665+0200 [INFO]  No output file was specified so the PR description will only be output to the console
2023-10-11T10:47:53.665+0200 [INFO]  Determining the current version of `hashicorp/go-azure-sdk` being used..
2023-10-11T10:47:53.665+0200 [INFO]  Old SDK Version is "v0.20230918.1115907"
2023-10-11T10:47:53.665+0200 [INFO]  Checking the changes between "v0.20230918.1115907" and "v0.20231005.1153009" of `hashicorp/go-azure-sdk`..
2023-10-11T10:48:14.806+0200 [INFO]  Updating `hashicorp/go-azure-sdk`..
2023-10-11T10:48:30.497+0200 [INFO]  Committed as "a0d657666ba969ca30154a9aec14f157a4a5e10d"
2023-10-11T10:48:33.434+0200 [INFO]  Processing Service "connectedvmware"..
2023-10-11T10:48:34.052+0200 [INFO]  Processed Service "connectedvmware".
2023-10-11T10:48:37.045+0200 [INFO]  Processing Service "cosmosdb"..
2023-10-11T10:48:37.621+0200 [INFO]  Attempting to update API Version "2022-05-15" to "2023-09-15" for Service "cosmosdb"..
2023-10-11T10:48:43.271+0200 [INFO]  Updated the Imports for Service "cosmosdb" to use API Version "2023-09-15" rather than "2022-05-15"..
2023-10-11T10:48:43.271+0200 [INFO]  Running `go mod tidy` / `go mod vendor`..
2023-10-11T10:48:54.889+0200 [INFO]  Running `make test` within "../../../"..
2023-10-11T10:51:45.702+0200 [INFO]  Committed as "cdb7749d478e961c10fb70e63f624f829d3dc639"
2023-10-11T10:51:45.702+0200 [INFO]  Updated Service "cosmosdb" from "2022-05-15" to "2023-09-15"
2023-10-11T10:51:45.702+0200 [INFO]  Attempting to update API Version "2022-11-15" to "2023-09-15" for Service "cosmosdb"..
2023-10-11T10:51:53.210+0200 [INFO]  Updated the Imports for Service "cosmosdb" to use API Version "2023-09-15" rather than "2022-11-15"..
2023-10-11T10:51:53.210+0200 [INFO]  Running `go mod tidy` / `go mod vendor`..
2023-10-11T10:52:18.365+0200 [INFO]  Running `make test` within "../../../"..
2023-10-11T10:53:38.033+0200 [INFO]  Resetting the working directory since `make test` failed..
2023-10-11T10:53:49.732+0200 [INFO]  Attempting to update API Version "2023-04-15" to "2023-09-15" for Service "cosmosdb"..
2023-10-11T10:53:55.432+0200 [INFO]  Updated the Imports for Service "cosmosdb" to use API Version "2023-09-15" rather than "2023-04-15"..
2023-10-11T10:53:55.432+0200 [INFO]  Running `go mod tidy` / `go mod vendor`..
2023-10-11T10:54:05.676+0200 [INFO]  Running `make test` within "../../../"..
2023-10-11T10:57:08.025+0200 [INFO]  Committed as "38ea9bc573d99a0a8ff512ae380ad7184c375f37"
2023-10-11T10:57:08.025+0200 [INFO]  Updated Service "cosmosdb" from "2023-04-15" to "2023-09-15"
2023-10-11T10:57:08.025+0200 [INFO]  Processed Service "cosmosdb".
2023-10-11T10:57:08.025+0200 [INFO]  Skipping writing PR description since an output file was not specified
2023-10-11T10:57:08.025+0200 [INFO]  Processing completed - summary of changes:
2023-10-11T10:57:08.025+0200 [INFO]  This PR updates the version of `hashicorp/go-azure-sdk` to `v0.20231005.1153009`.

This updates the following Services and API Versions:
* The service `cosmosdb` was updated to API Version `2023-09-15` (from `2022-05-15`).
* The service `cosmosdb` was updated to API Version `2023-09-15` (from `2023-04-15`).

## FAILED - API Versions
The following new API Versions are available but had compile-time errors when updating:
* The service `cosmosdb` - updating to API Version `2023-09-15` from `2022-11-15`.
[][][]
running `make test` in "../../../": stdout:
---
==> Checking that code complies with gofmt requirements...
==> Checking that Custom Timeouts are used...
==> Checking that acceptance test packages are used...
==> Checking for use of gradually deprecated functions...
==> Checking for use of deprecated functions...
==> Running Unit Tests...

---

stderr:
---
# github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos
internal/services/cosmos/cosmosdb_mongo_role_definition_resource.go:154:27: undefined: mongorbacs.MongoRoleDefinitionTypeOne
internal/services/cosmos/cosmosdb_mongo_role_definition_resource.go:161:30: cannot infer T (/Users/tharvey/code/src/github.com/hashicorp/terraform-provider-azurerm/vendor/github.com/hashicorp/go-azure-helpers/lang/pointer/generic.go:17:9)
internal/services/cosmos/cosmosdb_mongo_role_definition_resource.go:209:27: undefined: mongorbacs.MongoRoleDefinitionTypeOne
internal/services/cosmos/cosmosdb_mongo_role_definition_resource.go:216:30: cannot infer T (/Users/tharvey/code/src/github.com/hashicorp/terraform-provider-azurerm/vendor/github.com/hashicorp/go-azure-helpers/lang/pointer/generic.go:17:9)
make: *** [test] Error 1

---
[][][]
```

(note: `[][][]` in the output above is three backticks, replaced to not break the markdown layout.)

### Command Line Arguments

* `--azurerm-repo-path` - (Required) - Specifies the path to the root of the AzureRM Provider repository (typically this is `../../../`, when run from this repository). Example: `../../../`.
* `--go-sdk-repo-path` - (Optional) - Specifies the path to the root of the `hashicorp/go-azure-sdk` repository, if provided this path is used - if not a fresh copy is cloned into a temp directory. Example: `../../../../go-azure-sdk`.
* `--new-sdk-version` - (Required) - Specifies the version of `hashicorp/go-azure-sdk` that the Provider should be updated to. Example: `v0.20231005.1153009`.
* `--output-file` - (Optional) - Specifies the path to the output file containing the summary of changes performed by this tool, primarily intended to be used as a Pull Request body. Example: `pr-description.txt`.
