## Tool: `update-go-azure-sdk`

This tool updates the version of `hashicorp/go-azure-sdk` which is vendored in this repository, and then subsequently tries to update the Provider to account for any new API Versions.

Whilst this isn't going to be infallible, the intention is that when a new API Version is added to `hashicorp/go-azure-sdk` (via `hashicorp/pandora`) then the AzureRM Provider can be automatically updated to use that API Version - with any errors output for human review.

This tool makes a number of assumptions about the repository design, namely that an `internal` folder exists and that `hashicorp/go-azure-sdk` is not referenced outside of this directory - and in addition also tool assumes that Git is available and pre-configured.

Whilst this tool is intended to primarily be run in automation - [the `update-api-version` tool](../update-api-version) is better suited for interactive usage (since it won't perform rollbacks, and assumes the `hashicorp/go-azure-sdk` dependency is up-to-date).

### Example Usage

```bash
$ go build . && ./update-go-azure-sdk --new-sdk-version=v0.20231005.1153009 --azurerm-repo-path=../../../ --output-file=pr-description.txt
2023-10-10T15:14:12.249+0200 [INFO]  New SDK Version is "v0.20231005.1153009"
2023-10-10T15:14:12.249+0200 [INFO]  Output File Name is "pr-description.txt"
2023-10-10T15:14:12.249+0200 [INFO]  Working Directory is "../../../"
2023-10-10T15:14:12.249+0200 [INFO]  Old SDK Version is "v0.20230918.1115907"
2023-10-10T15:14:12.249+0200 [INFO]  Checking the changes between "v0.20230918.1115907" and "v0.20231005.1153009" of `hashicorp/go-azure-sdk`..
2023-10-10T15:14:34.275+0200 [INFO]  Updating `hashicorp/go-azure-sdk`..
2023-10-10T15:14:50.939+0200 [INFO]  Committed as "b5b40b0b49d9f192f9d65d620ecba98862764e6b"
2023-10-10T15:14:53.958+0200 [INFO]  Processing Service "connectedvmware"..
2023-10-10T15:14:54.552+0200 [INFO]  Processed Service "connectedvmware".
2023-10-10T15:14:57.578+0200 [INFO]  Processing Service "cosmosdb"..
2023-10-10T15:14:58.148+0200 [INFO]  Attempting to update API Version "2022-05-15" to "2023-09-15" for Service "cosmosdb"..
2023-10-10T15:15:03.823+0200 [INFO]  Updated the Imports for Service "cosmosdb" to use API Version "2023-09-15" rather than "2022-05-15"..
2023-10-10T15:15:03.823+0200 [INFO]  Running `go mod tidy` / `go mod vendor`..
2023-10-10T15:15:15.022+0200 [INFO]  Running `make test` within "../../../"..
2023-10-10T15:18:04.450+0200 [INFO]  Committed as "a88ab1b1be8344a07bcf316a62c3ea31a87be507"
2023-10-10T15:18:04.450+0200 [INFO]  Updated Service "cosmosdb" from "2022-05-15" to "2023-09-15"
2023-10-10T15:18:04.450+0200 [INFO]  Attempting to update API Version "2022-11-15" to "2023-09-15" for Service "cosmosdb"..
2023-10-10T15:18:10.035+0200 [INFO]  Updated the Imports for Service "cosmosdb" to use API Version "2023-09-15" rather than "2022-11-15"..
2023-10-10T15:18:10.035+0200 [INFO]  Running `go mod tidy` / `go mod vendor`..
2023-10-10T15:18:22.087+0200 [INFO]  Running `make test` within "../../../"..
2023-10-10T15:19:46.236+0200 [INFO]  Resetting the working directory since `make test` failed..
2023-10-10T15:19:56.549+0200 [INFO]  Attempting to update API Version "2023-04-15" to "2023-09-15" for Service "cosmosdb"..
2023-10-10T15:20:02.211+0200 [INFO]  Updated the Imports for Service "cosmosdb" to use API Version "2023-09-15" rather than "2023-04-15"..
2023-10-10T15:20:02.211+0200 [INFO]  Running `go mod tidy` / `go mod vendor`..
2023-10-10T15:20:12.584+0200 [INFO]  Running `make test` within "../../../"..
```

### Command Line Arguments

* `--azurerm-repo-path` - (Required) - Specifies the path to the root of the AzureRM Provider repository (typically this is `../../../`, when run from this repository). Example: `../../../`.
* `--new-sdk-version` - (Required) - Specifies the version of `hashicorp/go-azure-sdk` that the Provider should be updated to. Example: `v0.20231005.1153009`.
* `--output-file` - (Optional) - Specifies the path to the output file containing the summary of changes performed by this tool, primarily intended to be used as a Pull Request body. Example: `pr-description.txt`.
