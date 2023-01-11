# Debug

To debug the provider
## VSCODE

create debug config in VSCODE:


```json
        {
            "name": "Remote",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "${workspaceFolder}",
            "port": 8888,
            "host": "127.0.0.1"
        }
```

## Terminal 1

Compile the provider in debug mode

```sh
go build -gcflags="all=-N -l" ./cmd/debug
```

Launch `dlv` as:

```sh
dlv exec --listen=:8888 --headless --api-version=2 ./debug -- -debug
```

# VSCODE

Attach vs code to the debug session.

# Terminal 2

create a file `terraformrc` with the following content:
```
provider_installation {
  dev_overrides {               
    "snyk/azurerm" = "/Users/ricardol/git/terraform-provider-azurerm/debug" 
  }
}
```


```sh
export TF_CLI_CONFIG_FILE=/Users/ricardol/git/terraform-provider-azurerm/tmp/ex1/terraformrc
# In terminal 1 there's the value this variable should have
export TF_REATTACH_PROVIDERS=....
```


In a new folder, create a `provider.tf` with:
```
terraform {
  required_providers {
    azurerm = {
      source  = "snyk/azurerm"
      # version = "=3.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}
```

Create a `main.tf` with:
```
resource "azurerm_storage_account" "mysa" {
}
```

Run in the shell:
```
terraform import azurerm_storage_account.mysa "/subscriptions/ec86c868-4762-493c-ae65-411f88b2ffab/resourceGroups/rleal-resources/providers/Microsoft.Storage/storageAccounts/rlealstoraccount"
```
