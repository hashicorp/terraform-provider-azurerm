# Debug

To debug the provider
## VSCODE

create debug config:


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
export TF_REATTACH_PROVIDERS=....
```
