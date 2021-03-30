# Azure Identity Client Module for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/Azure/azure-sdk-for-go/sdk/azidentity)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity)
[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/go%20-%20azidentity?branchName=master)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=1846&branchName=master)
[![Code Coverage](https://img.shields.io/azure-devops/coverage/azure-sdk/public/1846/master)](https://img.shields.io/azure-devops/coverage/azure-sdk/public/1846/master)

The `azidentity` module provides a set of credential types for use with
Azure SDK clients that support Azure Active Directory (AAD) token authentication.

[Source code](https://github.com/Azure/azure-sdk-for-go/tree/master/sdk/azidentity)
| [Azure Active Directory documentation](https://docs.microsoft.com/azure/active-directory/)

# Getting started

## Install the package

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

Install the Azure Identity module:

```sh
go get -u github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

## Prerequisites

- an [Azure subscription](https://azure.microsoft.com/free/)
- Go 1.13 or above

### Authenticating during local development

When debugging and executing code locally it is typical for developers to use
their own accounts for authenticating calls to Azure services. The `azidentity`
module supports authenticating through developer tools to simplify
local development.

#### Authenticating via the Azure CLI

`DefaultAzureCredential` and `AzureCLICredential` can authenticate as the user
signed in to the [Azure CLI][azure_cli]. To sign in to the Azure CLI, run
`az login`. On a system with a default web browser, the Azure CLI will launch
the browser to authenticate a user.

![Azure CLI Account Sign In](img/AzureCLILogin.png)

When no default browser is available, `az login` will use the device code
authentication flow. This can also be selected manually by running `az login --use-device-code`.

![Azure CLI Account Device Code Sign In](img/AzureCLILoginDeviceCode.png)

## Key concepts

### Credentials

A credential is a type which contains or can obtain the data needed for a
service client to authenticate requests. Service clients across the Azure SDK
accept a credential instance when they are constructed, and use that credential
to authenticate requests.

The `azidentity` module focuses on OAuth authentication with Azure Active
Directory (AAD). It offers a variety of credential types capable of acquiring
an AAD access token. See [Credential Types](#credential-types "Credential Types") below for a list of this module's credential types.

### DefaultAzureCredential

`DefaultAzureCredential` is appropriate for most applications which will run in
the Azure Cloud because it combines common production credentials with
development credentials. `DefaultAzureCredential` attempts to authenticate via
the following mechanisms in this order, stopping when one succeeds:

![DefaultAzureCredential authentication flow](img/DAC_flow.PNG)

- Environment - `DefaultAzureCredential` will read account information specified
  via [environment variables](#environment-variables "environment variables")
  and use it to authenticate.
- Managed Identity - if the application is deployed to an Azure host with
  Managed Identity enabled, `DefaultAzureCredential` will authenticate with it.
- Azure CLI - If a user has signed in via the Azure CLI `az login` command,
  `DefaultAzureCredential` will authenticate as that user.

## Credential Types

### Authenticating Azure Hosted Applications

|credential|usage
|-|-
|DefaultAzureCredential|simplified authentication to get started developing applications for the Azure cloud
|ChainedTokenCredential|define custom authentication flows composing multiple credentials
|EnvironmentCredential|authenticate a service principal or user configured by environment variables
|ManagedIdentityCredential|authenticate the managed identity of an Azure resource

### Authenticating Service Principals

|credential|usage
|-|-
|ClientSecretCredential| authenticate a service principal using a secret
|CertificateCredential| authenticate a service principal using a certificate

### Authenticating Users

|credential|usage
|-|-
|InteractiveBrowserCredential|interactively authenticate a user with the default web browser
|DeviceCodeCredential| interactively authenticate a user on a device with limited UI
|UsernamePasswordCredential| authenticate a user with a username and password

### Authenticating via Development Tools

|credential|usage
|-|-
|AzureCLICredential|authenticate as the user signed in to the Azure CLI

## Environment Variables

DefaultAzureCredential] and 
EnvironmentCredential can be configured with
environment variables. Each type of authentication requires values for specific
variables:

#### Service principal with secret
|variable name|value
|-|-
|`AZURE_CLIENT_ID`|id of an Azure Active Directory application
|`AZURE_TENANT_ID`|id of the application's Azure Active Directory tenant
|`AZURE_CLIENT_SECRET`|one of the application's client secrets

#### Service principal with certificate
|variable name|value
|-|-
|`AZURE_CLIENT_ID`|id of an Azure Active Directory application
|`AZURE_TENANT_ID`|id of the application's Azure Active Directory tenant
|`AZURE_CLIENT_CERTIFICATE_PATH`|path to a PEM-encoded certificate file including private key (without password protection)

#### Username and password
|variable name|value
|-|-
|`AZURE_CLIENT_ID`|id of an Azure Active Directory application
|`AZURE_USERNAME`|a username (usually an email address)
|`AZURE_PASSWORD`|that user's password

Configuration is attempted in the above order. For example, if values for a
client secret and certificate are both present, the client secret will be used.

## Troubleshooting

### Error Handling

Credentials return `CredentialUnavailableError` when they're unable to attempt
authentication because they lack required data or state. For example,
`NewEnvironmentCredential` will return an error when
[its configuration](#environment-variables "its configuration") is incomplete.

Credentials return `AuthenticationFailedError` when they fail
to authenticate. Call `Error()` on `AuthenticationFailedError` to see why authentication failed. When returned by
`DefaultAzureCredential` or `ChainedTokenCredential`,
the message collects error messages from each credential in the chain.

For more details on handling specific Azure Active Directory errors please refer to the
Azure Active Directory
[error code documentation](https://docs.microsoft.com/azure/active-directory/develop/reference-aadsts-error-codes).

### Logging

This module uses the classification based logging implementation in azcore. To turn on logging set `AZURE_SDK_GO_LOGGING` to `all`. If you only want to include logs for `azidentity`, you must create you own logger and set the log classification as `LogCredential`.
Credentials log basic information only, including `GetToken` success or failure and errors. These log entries do not contain authentication secrets but may contain sensitive information.

To obtain more detailed logging, including request/response bodies and header values, make sure to leave the logger as default or enable the `LogRequest` and/or `LogResponse` classificatons. A logger that only includes credential logs can be like the following:

```go
// Set log to output to the console
azcore.Log().SetListener(func(cls LogClassification, s string) {
		fmt.Println(s) // printing log out to the console
  })
  
// Include only azidentity credential logs
azcore.Log().SetClassifications(azidentity.LogCredential)
```

> CAUTION: logs from credentials contain sensitive information.
> These logs must be protected to avoid compromising account security.

# Next steps

## Provide Feedback

If you encounter bugs or have suggestions, please
[open an issue](https://github.com/Azure/azure-sdk-for-go/issues) and assign the `Azure.Identity` label.

# Contributing

This project welcomes contributions and suggestions. Most contributions require
you to agree to a Contributor License Agreement (CLA) declaring that you have
the right to, and actually do, grant us the rights to use your contribution.
For details, visit [https://cla.microsoft.com](https://cla.microsoft.com).

When you submit a pull request, a CLA-bot will automatically determine whether
you need to provide a CLA and decorate the PR appropriately (e.g., label,
comment). Simply follow the instructions provided by the bot. You will only
need to do this once across all repos using our CLA.

This project has adopted the
[Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information, see the
[Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/)
or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any
additional questions or comments.

[azure_cli]: https://docs.microsoft.com/cli/azure
[azblob]: https://github.com/Azure/azure-sdk-for-go/tree/master/sdk

![Impressions](https://azure-sdk-impressions.azurewebsites.net/api/impressions/azure-sdk-for-go%2Fsdk%2Fidentity%2Fazure-identity%2FREADME.png)
