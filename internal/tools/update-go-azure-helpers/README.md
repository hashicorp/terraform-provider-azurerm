## Tool: `update-go-azure-helpers`

This tool updates the version of `hashicorp/go-azure-helpers` which is vendored in this repository.

This tool makes a number of assumptions about the repository design, namely that an `internal` folder exists and that `hashicorp/go-azure-helpers` is not referenced outside of this directory - and in addition also tool assumes that Git is available and pre-configured.

### Example Usage

```bash
 $ go build . && ./update-go-azure-sdk --new-helpers-version=v0.66.3 --azurerm-repo-path=../../../ --azure-helpers-repo-path=../../../../go-azure-helpers
2024-02-26T14:07:26.430+0100 [INFO]  New Go Azure Helpers Version is "v0.66.3"
2024-02-26T14:07:26.430+0100 [INFO]  The `hashicorp/terraform-provider-azurerm` repository is located at "../../../"
2024-02-26T14:07:26.430+0100 [INFO]  The `hashicorp/go-azure-helpers` repository is located at "../../../../go-azure-helpers"
2024-02-26T14:07:26.430+0100 [INFO]  Determining the current version of `hashicorp/go-azure-helpers` being used..
2024-02-26T14:07:26.431+0100 [INFO]  Old Go Azure Helpers Version is "v0.66.2"
2024-02-26T14:07:26.431+0100 [INFO]  Updating `hashicorp/go-azure-helpers`..
2024-02-26T14:07:36.417+0100 [INFO]  Committed as "7fd5d3a5eac63b1660690ecb9579b88340ac776b"
```

### Command Line Arguments

* `--azurerm-repo-path` - (Required) - Specifies the path to the root of the AzureRM Provider repository (typically this is `../../../`, when run from this repository). Example: `../../../`.
* `--azure-helpers-repo-path` - (Optional) - Specifies the path to the root of the `hashicorp/go-azure-helpers` repository, if provided this path is used - if not a fresh copy is cloned into a temp directory. Example: `../../../../go-azure-helpers`.
* `--new-helpers-version` - (Required) - Specifies the version of `hashicorp/go-azure-helpers` that the Provider should be updated to. Example: `v0.20231005.1153009`.
