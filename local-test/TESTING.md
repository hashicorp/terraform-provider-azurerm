# Local Testing Guide for definitionVersion Feature

## Prerequisites
- Go 1.25.x installed
- Azure subscription with appropriate permissions
- Azure CLI configured (`az login`)
- Terraform installed

## Option 1: Use Development Override (Easiest)

This allows Terraform to use your locally built provider without installation.

### Steps:

1. **Build the provider:**
   ```powershell
   go build -o terraform-provider-azurerm.exe
   ```

2. **Create a dev overrides configuration:**
   
   Create/edit `%APPDATA%\terraform.rc` (Windows) or `~/.terraformrc` (Linux/Mac):
   
   ```hcl
   provider_installation {
     dev_overrides {
       "hashicorp/azurerm" = "C:\\gitRepos\\go\\terraform-provider-azurerm"
     }
     
     # For all other providers, install them normally
     direct {
       exclude = ["hashicorp/azurerm"]
     }
   }
   ```
   
   Update the path to match where your compiled binary is located.

3. **Test in the local-test directory:**
   ```powershell
   cd local-test
   terraform init
   terraform plan
   terraform apply
   ```

4. **Verify the feature:**
   - Check that `definition_version` appears in the plan
   - Apply and verify in Azure Portal
   - Check the output shows the version
   - Test data source reads the version correctly

## Option 2: Install to Local Plugin Directory

### Steps:

1. **Build for your platform:**
   ```powershell
   # Windows
   $env:GOOS="windows"
   $env:GOARCH="amd64"
   go build -o terraform-provider-azurerm_v99.99.99.exe
   
   # Or Linux
   $env:GOOS="linux"
   $env:GOARCH="amd64"
   go build -o terraform-provider-azurerm_v99.99.99
   ```

2. **Create plugin directory structure:**
   ```powershell
   # Windows
   $pluginDir = "$env:APPDATA\terraform.d\plugins\registry.terraform.io\hashicorp\azurerm\99.99.99\windows_amd64"
   New-Item -ItemType Directory -Force -Path $pluginDir
   Copy-Item terraform-provider-azurerm_v99.99.99.exe "$pluginDir\terraform-provider-azurerm_v99.99.99.exe"
   ```

3. **Update terraform.tf to use local version:**
   ```hcl
   terraform {
     required_providers {
       azurerm = {
         source  = "hashicorp/azurerm"
         version = "99.99.99"
       }
     }
   }
   ```

## Option 3: Run Unit Tests

Test without Azure credentials:

```powershell
# Test specific files
go test ./internal/services/policy/...

# Run with verbose output
go test -v ./internal/services/policy/assignment_resource_base_test.go

# Test specific function
go test -v -run TestAccSubscriptionPolicyAssignment_definitionVersion ./internal/services/policy/...
```

## Option 4: Run Acceptance Tests (Requires Azure)

**Note:** These tests create real Azure resources and may incur costs.

```powershell
# Set required environment variables
$env:ARM_SUBSCRIPTION_ID = "your-subscription-id"
$env:ARM_CLIENT_ID = "your-client-id"
$env:ARM_CLIENT_SECRET = "your-client-secret"
$env:ARM_TENANT_ID = "your-tenant-id"
$env:TF_ACC = "1"

# Run specific tests
go test -v -timeout 120m ./internal/services/policy/ -run TestAccSubscriptionPolicyAssignment_definitionVersion

# Or use make (if available)
make acctests SERVICE='policy' TESTARGS='-run=TestAccSubscriptionPolicyAssignment_definitionVersion'
```

## Testing Checklist

Test these scenarios in the local-test directory:

### Basic Functionality
- [ ] Create assignment with `definition_version = "1.0.0"` (exact version)
- [ ] Create assignment with `definition_version = "1.0.*"` (patch wildcard)
- [ ] Create assignment with `definition_version = "1.*.*"` (minor wildcard)
- [ ] Create assignment without `definition_version` (should work - field is optional)
- [ ] Update existing assignment to add `definition_version`
- [ ] Update `definition_version` from one value to another
- [ ] Remove `definition_version` from existing assignment

### Data Source
- [ ] Read assignment with `definition_version` using data source
- [ ] Verify output shows correct version

### Edge Cases
- [ ] Invalid version format (should fail validation)
- [ ] Import existing assignment (with and without version)

### All Resource Types
Test with:
- [ ] `azurerm_subscription_policy_assignment`
- [ ] `azurerm_management_group_policy_assignment`
- [ ] `azurerm_resource_group_policy_assignment`
- [ ] `azurerm_resource_policy_assignment`

## Verification

After apply, verify in Azure Portal or CLI:

```bash
# Using Azure CLI
az policy assignment show --name test-def-version --scope /subscriptions/SUBSCRIPTION_ID

# Look for the definitionVersion field in the output
```

## Cleanup

```powershell
# Remove terraform resources
cd local-test
terraform destroy

# Remove dev override (edit terraform.rc and remove the dev_overrides block)

# Clean up local plugin directory if you used Option 2
Remove-Item -Recurse "$env:APPDATA\terraform.d\plugins\registry.terraform.io\hashicorp\azurerm\99.99.99"
```

## Troubleshooting

### Provider not found
- Verify the path in terraform.rc matches your binary location
- Check that the binary was built successfully
- Ensure you're using absolute paths in terraform.rc

### Definition version not working
- Verify you're using your local provider (Terraform will show a warning about dev overrides)
- Check Azure API permissions
- Ensure the policy definition has versioning metadata

### Tests timing out
- Increase timeout: `go test -timeout 180m ...`
- Run individual tests instead of the full suite
- Check Azure credentials are valid
