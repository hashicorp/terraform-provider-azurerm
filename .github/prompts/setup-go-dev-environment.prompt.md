---
mode: agent
tools: [runCommands]
description: This document provides instructions for installing Go and setting up the development environment required for the Terraform AzureRM provider project. It includes steps to install Go, set up GOPATH, and install required development tools.
---

# Install Go Development Environment

To set up the development environment for the Terraform AzureRM provider, follow these steps:

## 1. Install Go

### On Windows:
- Download Go 1.22.x or later from https://golang.org/dl/
- Run the installer and follow the installation wizard
- Or use Chocolatey:
```powershell
choco install golang terraform make -y
refreshenv
```

### On macOS:
```bash
# Using Homebrew
brew install go terraform

# Or download from https://golang.org/dl/
```

### On Linux:
```bash
# Download and install Go
wget https://golang.org/dl/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Add to PATH in ~/.bashrc or ~/.zshrc
export PATH=$PATH:/usr/local/go/bin
```

## 2. Set up GOPATH (if using older Go workspace model)
```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

## 3. Clone the Repository
```bash
# Using modern Go modules (recommended)
git clone https://github.com/hashicorp/terraform-provider-azurerm.git
cd terraform-provider-azurerm

# Or using traditional GOPATH
mkdir -p $GOPATH/src/github.com/hashicorp
cd $GOPATH/src/github.com/hashicorp
git clone https://github.com/hashicorp/terraform-provider-azurerm.git
cd terraform-provider-azurerm
```

## 4. Install Development Tools
```bash
# Install required tooling
make tools

# This installs tools like:
# - gofmt, goimports (code formatting)
# - golangci-lint (linting)
# - terraform-plugin-docs (documentation generation)
# - Other development dependencies
```

## 5. Verify Installation
```bash
# Check Go version
go version

# Check Terraform version
terraform version

# Build the provider to verify setup
make build

# Run tests to ensure everything works
make test
```

## 6. Additional Windows Setup
If you're on Windows, you'll also need:
- Git Bash for Windows
- Make for Windows (or use the Chocolatey installation above)
- Set git config: `git config --global core.autocrlf false`
- Set git config: `git config --system core.longpaths true`

## 7. Environment Variables for Testing
Set up the following environment variables for running acceptance tests:
```bash
export ARM_CLIENT_ID="your-service-principal-client-id"
export ARM_CLIENT_SECRET="your-service-principal-client-secret"
export ARM_SUBSCRIPTION_ID="your-azure-subscription-id"
export ARM_TENANT_ID="your-azure-tenant-id"
export ARM_TEST_LOCATION=WestEurope
export ARM_TEST_LOCATION_ALT=EastUS2
export ARM_TEST_LOCATION_ALT2=WestUS2
```
