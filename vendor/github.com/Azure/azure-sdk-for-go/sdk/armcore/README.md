# Azure Resource Manager Core Client Module for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/Azure/azure-sdk-for-go/sdk/armcore)](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/armcore)
[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/go%20-%20armcore%20-%20ci?branchName=master)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=1844&branchName=master)
[![Code Coverage](https://img.shields.io/azure-devops/coverage/azure-sdk/public/1844/master)](https://img.shields.io/azure-devops/coverage/azure-sdk/public/1844/master)

The `armcore` module provides functions and types for Go SDK ARM client modules.
These modules follow the [Azure SDK Design Guidelines for Go](https://azure.github.io/azure-sdk/golang_introduction.html).

## Getting started

This project uses [Go modules](https://github.com/golang/go/wiki/Modules) for versioning and dependency management.

Typically, you will not need to explicitly install `armcore` as it will be installed as an ARM client module dependency.
To add the latest version to your `go.mod` file, execute the following command.

```bash
go get -u github.com/Azure/azure-sdk-for-go/sdk/armcore
```

General documentation and examples can be found on [pkg.go.dev](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/armcore).

## Contributing
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
