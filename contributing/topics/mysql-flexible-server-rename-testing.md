# MySQL Flexible Server rename: testing, PR, and release runbook

Prepared against the local rename implementation on **2026-10-17**, with offline validation updated on **2026-10-19**. This is not evidence that live Azure tests, a PR, or a release have happened.

| Goal | Sections |
| --- | --- |
| Verify code without Azure changes | 3-5 |
| Exercise a local binary against real Azure | 6-7 |
| Run automated live acceptance | 8 |
| Test existing resources, dependencies, and recovery | 9-11 |
| Prepare the PR and publish/roll out the feature | 12-14 |
| Decide whether testing/release is complete | 15 |

## 1. Scope and current status

The change adds `rename_enabled` to `azurerm_mysql_flexible_server`. With that option enabled, changing `name` at the **same Terraform resource address** renames the existing server. Without the option, a name change continues to require replacement.

The provider sends a PUT to the destination resource ID using the ARM-registered non-preview API version `2025-12-30`, with only:

```json
{
  "location": "<location returned by the source GET>",
  "properties": {
    "createMode": "Rename",
    "sourceServerResourceId": "<canonical source resource ID>"
  }
}
```

It waits for the long-running operation and verifies the destination before changing the Terraform resource ID, resource identity, name, and FQDN. Ordinary server operations still use `2023-12-30`. The CLI's submission-only behavior is intentionally **not** the Terraform completion contract.

| Item | Status as of 2026-10-19 |
| --- | --- |
| Implementation and resource documentation | Present as local, uncommitted changes |
| Focused offline rename tests, real Terraform CLI with mocked HTTP, and adapter guard tests | Passed with Go 1.26.9 and Terraform 1.13.5 |
| Existing MySQL name/configuration validators and ID parsers | Passed |
| MySQL package and its new acceptance test | Compiled |
| Provider schema and production/acceptance factory schema checks | Passed |
| Live rename acceptance test | Not run; explicitly opted out |
| Standalone provider executable, loaded schema, and standalone configuration validation | Passed |
| Real Terraform CLI end-to-end flow against Azure | Not run |
| HA, private networking, CMK, and replication live scenarios | Not yet validated |
| Full repository lint, generation, and PR checks | Not yet run for this change |
| Feature PR, upstream approval, package publication, and backend rollout | Not performed |

The implementation-phase base commit was `bb0454fbdfeadd1ae7d1c67c6eb35c58b4995e9b`. That commit **does not include the uncommitted feature**. Record the actual feature revision and diff for every subsequent run.

The real CLI regression exposed a gap in the SDK-only tests: an unknown `id`
caused the Plugin SDK to request replacement. Planning a known destination ID
instead was also insufficient: the SDK kept the source ID and Terraform outputs
became stale. The rename-specific protocol adapter now preserves unknown ID/FQDN
values while removing only the synthetic ID replacement for a valid, opted-in
rename. Both provider factories use it; genuine replacement paths are preserved.
No vendored dependency changes were made. The offline CLI test seeds the source
locally, so it does not verify actual Azure creation, rename enablement, or data
persistence.

### Files in scope

| File relative to the repository root | Purpose |
| --- | --- |
| `internal\services\mysql\client\client.go` | Dedicated `2025-12-30` rename client, normal authentication/environment configuration |
| `internal\services\mysql\mysql_flexible_server_resource.go` | Opt-in schema, planning, mutable identity, update dispatch, and missing-source safeguard |
| `internal\services\mysql\mysql_flexible_server_rename.go` | Source/destination checks, name availability, PUT, polling, and state reconciliation |
| `internal\services\mysql\mysql_flexible_server_rename_test.go` | Offline planning and mocked-HTTP lifecycle/failure tests |
| `internal\services\mysql\mysql_flexible_server_rename_cli_test.go` | Real Terraform CLI with mocked HTTP: update-only plans, identity, import, and repeated rename |
| `internal\services\mysql\mysql_flexible_server_rename_server.go` | Narrow protocol adapter removing only the SDK's synthetic ID replacement on opted-in rename |
| `internal\services\mysql\mysql_flexible_server_rename_server_test.go` | Adapter scope, replacement preservation, and error propagation tests |
| `internal\provider\framework\factory_builder.go` | Adapter wiring for both production and acceptance-test provider factories |
| `internal\services\mysql\mysql_flexible_server_rename_resource_test.go` | Explicitly enabled live acceptance test |
| `website\docs\r\mysql_flexible_server.html.markdown` | User-facing behavior, limitations, and recovery |
| `contributing\topics\mysql-flexible-server-rename-testing.md` | This runbook |

### What "deployment" means

| Path | What is installed or changed | Requires a published provider? |
| --- | --- | --- |
| Offline Go tests | Local code with mocked HTTP; no Azure changes | No |
| Local provider test | Development binary runs on your workstation and calls real Azure | No |
| Live Go acceptance test | Test harness serves the provider built from the checkout and creates real Azure resources | No |
| AzureRM release | Maintainers build, sign, and publish provider packages | Yes, for normal customer installation |
| MySQL service rollout/enablement | Service owners deploy or enable the backend capability | Independent of provider publication |

There is no Terraform provider binary to deploy into the MySQL server. Publishing the Azure CLI command neither ships AzureRM support nor enables the backend feature.

## 2. Safety and approval gates

Complete these before any live section:

1. Use an approved disposable subscription, resource group, and server. Agree on ownership, cost, test window, and cleanup.
2. Confirm that the chosen region and subscription support the actual rename PUT using `2025-12-30`. An advertised API version or successful GET does not prove rename support.
3. Have the service owner confirm topology-specific support. Existing canary evidence using `2025-03-01-preview` is not proof that the provider's `2025-12-30` path works.
4. Confirm existing identity permissions to read the source and destination, check name availability, write the destination, and read the operation result. Creation and cleanup require their own permissions.
5. Do not grant new roles, change feature switches, open firewall access, or disable security controls just to make a test pass without separate approval.
6. Do not use a customer server, shared production Terraform state, or a server managed by another Terraform resource address.
7. Keep state, saved plans, tfvars, authentication material, SQL evidence, and provider logs in approved private storage outside this Git checkout. Terraform state and plan files can contain passwords even when values are marked sensitive.
8. Do not run a broad acceptance suite by accident. The live test creates, renames, and deletes real resources.
9. Plan the application endpoint change and any downtime. Do not assume the old hostname remains usable.
10. Resolve the separate dependency migration described in section 10 before testing a server with Terraform-managed children or replicas.

If the backend is not enabled, finish the offline steps and stop before live mutations. Record the specific unfulfilled prerequisite rather than claiming a successful feature test.

## 3. Prepare the Windows development shell

Use **PowerShell 7** for the PowerShell examples. Run sections 3-9 in the same shell unless directed otherwise. Git Bash/WSL commands are identified separately.

### 3.1 Inspect the checkout and preserve existing work

```powershell
$ErrorActionPreference = 'Stop'
$Repo = 'C:\Work\repos\terraform-provider-azurerm'
Set-Location $Repo

git status --short
git branch --show-current
git rev-parse HEAD
git remote -v
git diff --check
if ($LASTEXITCODE -ne 0) { throw 'Resolve whitespace errors before continuing.' }
```

Inspect the six feature source/test/resource-documentation files listed above, including untracked files: `git diff` alone does not include new files. Preserve unrelated work. Do not reset, clean, stash, or switch branches as a troubleshooting shortcut.

### 3.2 Select the required tools

```powershell
$RequiredGo = (Get-Content (Join-Path $Repo '.go-version') -Raw).Trim()
Write-Host "Required Go version: $RequiredGo"
Get-Command go, terraform, git, gh -ErrorAction SilentlyContinue
```

The inspected checkout requires **Go 1.26.9**, but `.go-version` remains authoritative.

If Go is already installed outside PATH, prepend its `bin` directory for this shell only:

```powershell
$GoBin = '<absolute-path-to-existing-go-bin>'
$env:PATH = "$GoBin;$env:PATH"
go version
if ($LASTEXITCODE -ne 0) { throw 'Go is not available.' }
$env:GOTOOLCHAIN = 'local'
```

If the required toolchain is absent, obtain the matching OS/architecture archive from [official Go downloads](https://go.dev/dl/), verify its published SHA-256, and extract it to a dedicated local tool directory. Do not change `go.mod` or `.go-version` to accommodate another installed version.

Install Terraform only if absent, using your approved installation method. One Windows option is:

```powershell
winget install Hashicorp.Terraform
```

Use a currently supported Terraform 1.x version, **at least 1.12 for resource-identity verification**. Reopen the shell after an installation if PATH needs refreshing, then repeat the environment setup. Also install GitHub CLI for PR work and Azure CLI/MySQL client only for the live checks that need them.

```powershell
go version
terraform version
git --version
gh --version
```

For the repository's complete Make targets, follow [Building the Provider](building-the-provider.md): they invoke POSIX shell utilities and are not native PowerShell commands. Use a supported Git Bash or Linux/WSL environment; some pinned tooling targets require Linux/macOS. Keep any tooling changes scoped to this task rather than modifying global Git settings.

### 3.3 Create a private test workspace

```powershell
$RunRoot = Join-Path $env:LOCALAPPDATA ("AzureRM-Rename-Test\" + [guid]::NewGuid())
$Bin = Join-Path $RunRoot 'bin'
$Work = Join-Path $RunRoot 'terraform'
$Evidence = Join-Path $RunRoot 'evidence'
New-Item -ItemType Directory -Path $Bin, $Work, $Evidence | Out-Null
Write-Host "Private test workspace: $RunRoot"

git rev-parse HEAD | Set-Content (Join-Path $Evidence 'base-commit.txt')
git status --short | Set-Content (Join-Path $Evidence 'working-tree.txt')
```

Apply your organization's access-control and retention requirements to this directory before writing secrets or state. Do not attach its entire contents to a public PR.

## 4. Run offline checks

### 4.1 Explicitly disable live acceptance testing

```powershell
Set-Location $Repo
Remove-Item Env:\TF_ACC -ErrorAction SilentlyContinue
Remove-Item Env:\ARM_MYSQL_RENAME_ENABLED -ErrorAction SilentlyContinue
Remove-Item Env:\GO_AZURE_SDK_SKIP_POLLING_DELAY -ErrorAction SilentlyContinue
$env:TF_ACC_TERRAFORM_PATH = (Get-Command terraform).Source
```

The offline tests set their own mocked polling behavior. Never carry a polling-delay override into a live run.

### 4.2 Run the exact focused suite

```powershell
go test -p 6 .\internal\services\mysql `
  -run '^TestMysqlFlexibleServerRename|^TestAccMySqlFlexibleServer_rename$' `
  -count=1 -timeout 10m
if ($LASTEXITCODE -ne 0) { throw 'Rename tests failed.' }

go test -p 2 .\internal\services\mysql\validate .\internal\services\mysql\parse -count=1
if ($LASTEXITCODE -ne 0) { throw 'MySQL validation/parser regression tests failed.' }
```

The original SDK-only suite and parser/validator checks **passed during the implementation session on 2026-10-17**. Rerun them on the revision being reviewed; prior results do not certify later edits. The first command now also runs a real Terraform CLI against mocked HTTP, with a locally seeded source, and skips the opt-in live test. Set `TF_ACC_TERRAFORM_PATH` to Terraform 1.12 or newer; that path does not enable Azure acceptance testing.

To record individual cases and the live skip explicitly:

```powershell
go test -p 6 .\internal\services\mysql `
  -run '^TestMysqlFlexibleServerRename|^TestAccMySqlFlexibleServer_rename$' `
  -count=1 -timeout 10m -v *> (Join-Path $Evidence 'offline-rename.txt')
$TestExit = $LASTEXITCODE
Get-Content (Join-Path $Evidence 'offline-rename.txt') -Tail 30
if ($TestExit -ne 0) { throw 'See the private offline test log.' }
```

Expected coverage:

| Area | Expected assertion |
| --- | --- |
| Default compatibility | Name changes still require replacement without opt-in |
| Opt-in planning | In-place update through the actual Terraform protocol; ID and FQDN unknown until apply; outputs follow the destination after apply |
| Plan rejection | Case-only rename, unknown destination name or opt-in, and simultaneous property changes rejected |
| Submission | Exact destination path, `2025-12-30` API, minimal payload, canonical source ID |
| Completion | Initial 200/201/202 responses with an operation URL are polled |
| Convergence | A-to-B-to-C renames converge without another rename |
| Follow-up operations | Update and delete use the final destination ID |
| Preconditions | Missing/not-ready source, wrong source ID, destination collision, and unavailable name fail |
| Availability | Source location used; specific `AuthorizationFailed` fallback covered |
| Failures | Rejection, failed/canceled LRO, polling authorization error, interruption, missing operation URL, and wrong destination fail |
| State on failure | Original ID, identity, name, and FQDN retained before destination verification |
| Refresh safety | Missing rename-enabled server reports an error instead of silently disappearing from state |
| Toggle-only change | Enabling/disabling rename alone reads but does not modify Azure |

Mocked "feature disabled" and similar responses verify error handling, not a real deployment's error codes or rollback behavior.

### 4.3 Check formatting and wider repository gates

```powershell
$GoFiles = @(
  'internal\services\mysql\client\client.go',
  'internal\services\mysql\mysql_flexible_server_resource.go',
  'internal\services\mysql\mysql_flexible_server_rename.go',
  'internal\services\mysql\mysql_flexible_server_rename_test.go',
  'internal\services\mysql\mysql_flexible_server_rename_resource_test.go'
)
$Unformatted = @(gofmt -l $GoFiles)
if ($LASTEXITCODE -ne 0 -or $Unformatted.Count -ne 0) {
  $Unformatted
  throw 'Format the listed Go files and rerun tests.'
}
git diff --check
if ($LASTEXITCODE -ne 0) { throw 'Whitespace check failed.' }
```

Before submitting the feature PR, run the documented full gates in the supported POSIX shell from the repository root:

```bash
unset TF_ACC ARM_MYSQL_RENAME_ENABLED GO_AZURE_SDK_SKIP_POLLING_DELAY
make pr-check
make markdownlint
make linkcheck-local
git diff --check
```

`make pr-check` includes generation, build, tests, linting, and website checks. It can take substantial time and generate changes. Inspect all resulting diffs and retain only intended updates; do not discard somebody else's work. The Make targets acquire their pinned tools from `.tools` as needed.

Do not report `make pr-check` as passed based only on the focused Go suite. If the local environment cannot run a gate, record that limitation and obtain the corresponding CI result before merge.

## 5. Build and prove that Terraform loads the changed provider

### 5.1 Build an explicit local executable

```powershell
Set-Location $Repo
$ProviderExe = Join-Path $Bin 'terraform-provider-azurerm.exe'
go build -p 6 -o $ProviderExe .
if ($LASTEXITCODE -ne 0) { throw 'Provider build failed.' }
Get-FileHash -Algorithm SHA256 -LiteralPath $ProviderExe |
  Format-List | Out-File (Join-Path $Evidence 'provider-sha256.txt')
go version -m $ProviderExe |
  Set-Content (Join-Path $Evidence 'provider-build-info.txt')
```

This explicit build compiles the provider without running all Make checks. It does not substitute for section 4.3.

### 5.2 Write a shell-scoped development override

```powershell
$DevConfig = Join-Path $RunRoot 'development.tfrc'
$QuotedBin = ConvertTo-Json -InputObject $Bin -Compress
@"
provider_installation {
  dev_overrides {
    "hashicorp/azurerm" = $QuotedBin
  }
  direct {}
}
"@ | Set-Content -Encoding utf8 -LiteralPath $DevConfig
$env:TF_CLI_CONFIG_FILE = $DevConfig
```

`ConvertTo-Json` supplies a quoted, escaped Windows path usable by HCL. Do not edit the global Terraform CLI configuration. Development overrides bypass normal package version/checksum selection for this provider; the binary's recorded hash is therefore important.

### 5.3 Create the minimal standalone configuration

Create `$Work\main.tf` with the following contents. This configuration creates a new disposable server; section 9 explains the different path for an existing server.

```hcl
terraform {
  required_version = ">= 1.12.0, < 2.0.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.2.0, < 5.0.0"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id                 = var.subscription_id
  resource_provider_registrations = "none"
}

variable "subscription_id" {
  type = string
}

variable "resource_group_name" {
  type = string
}

variable "location" {
  type = string
}

variable "server_name" {
  type = string
}

variable "administrator_password" {
  type      = string
  sensitive = true
}

variable "rename_enabled" {
  type    = bool
  default = false
}

resource "azurerm_resource_group" "test" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_mysql_flexible_server" "test" {
  name                   = var.server_name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "renameadmin"
  administrator_password = var.administrator_password
  sku_name               = "B_Standard_B1ms"
  rename_enabled         = var.rename_enabled

  lifecycle {
    prevent_destroy = true
  }

  timeouts {
    create = "2h"
    update = "2h"
    delete = "1h"
  }
}

output "server_id" {
  value = azurerm_mysql_flexible_server.test.id
}

output "server_name" {
  value = azurerm_mysql_flexible_server.test.name
}

output "server_fqdn" {
  value = azurerm_mysql_flexible_server.test.fqdn
}
```

The version constraint is for the existing AzureRM 4.x development baseline; it does **not** identify a released version containing rename. The override is what supplies the feature. Replace the constraint with the actual tested release in section 14.

Provider registration is deliberately disabled to avoid an incidental subscription change. Confirm the required providers are already registered through your approved process. The burstable SKU is for the standalone test only; replication and HA require suitable supported SKUs. Do not change the SKU during a rename.

For this module-free, local-state, single-provider configuration, a development override allows direct schema/validate/plan operations; `terraform init` is not required for the local provider. If you add modules, other providers, or a remote backend, initialize those deliberately. `terraform init` still resolves a published provider and is not proof that the development binary will be used by subsequent commands.

### 5.4 Verify the schema before any Azure mutation

```powershell
Set-Location $Work
terraform fmt -check
if ($LASTEXITCODE -ne 0) { throw 'Format main.tf before proceeding.' }
$SchemaText = terraform providers schema -json
if ($LASTEXITCODE -ne 0) { throw 'Cannot load development provider schema.' }
$Schema = ($SchemaText -join "`n") | ConvertFrom-Json
$ServerSchema = $Schema.provider_schemas.'registry.terraform.io/hashicorp/azurerm'.
  resource_schemas.azurerm_mysql_flexible_server
if (-not $ServerSchema.block.attributes.rename_enabled.optional) {
  throw 'Wrong provider loaded: rename_enabled is absent or not optional.'
}
terraform validate
if ($LASTEXITCODE -ne 0) { throw 'Configuration validation failed.' }
```

Expect Terraform's development-override warning. If the new argument is missing, check `$env:TF_CLI_CONFIG_FILE`, the executable directory, and its hash. A version string alone cannot prove the modified binary was loaded.

## 6. Authenticate and confirm live-service readiness

This section starts the live path. Stop unless section 2 is approved.

### 6.1 Select one intended authentication method

For an interactive manual test, use the existing Azure CLI login in the intended cloud/tenant and remove stale service-principal/OIDC/MSI overrides from this dedicated shell. Do not clear or alter credentials used by other shells.

```powershell
az cloud show --query name --output tsv
az login --tenant '<approved-tenant-id>'
if ($LASTEXITCODE -ne 0) { throw 'Azure login failed.' }
az account set --subscription '<approved-subscription-id>'
if ($LASTEXITCODE -ne 0) { throw 'Cannot select test subscription.' }
az account show --query '{subscription:id,tenant:tenantId,name:name}' --output json
if ($LASTEXITCODE -ne 0) { throw 'Cannot verify selected account.' }

$env:ARM_SUBSCRIPTION_ID = '<approved-subscription-id>'
$env:ARM_TENANT_ID = '<approved-tenant-id>'
$env:ARM_USE_CLI = 'true'
$env:TF_VAR_subscription_id = $env:ARM_SUBSCRIPTION_ID
$env:TF_VAR_administrator_password = Read-Host 'Disposable server administrator password' -MaskInput
```

For sovereign clouds, use their approved login/environment configuration and ARM endpoint. Do not switch clouds or tenants implicitly.

For unattended/acceptance testing, use a separately approved service principal or supported workload identity as described in section 8. Do not create credentials as part of this runbook without approval.

### 6.2 Fill in non-secret test inputs

```powershell
$Suffix = Get-Date -Format 'yyyyMMddHHmmss'
$SourceName = "tf-rename-$Suffix"
$DestinationName = "$SourceName-b"
$FinalName = "$SourceName-c"
$ResourceGroup = "tf-rename-$Suffix"
$Location = '<approved-rename-enabled-region>'

@"
resource_group_name = "$ResourceGroup"
location            = "$Location"
server_name         = "$SourceName"
rename_enabled      = false
"@ | Set-Content -Encoding utf8 -LiteralPath (Join-Path $Work 'terraform.tfvars')
```

Names must be globally available, lowercase, 3-63 characters, and meet the existing server-name validator. Ensure the test password meets the resource requirements; do not write it into tfvars or command arguments.

Before create, have the service owner confirm feature enablement, region/SKU capacity, quota, expected API routing, and supported topology. Neither this configuration nor the new provider option enables a service feature switch.

## 7. Manual end-to-end Terraform test

### 7.1 Create the disposable source

```powershell
Set-Location $Work
terraform validate
if ($LASTEXITCODE -ne 0) { throw 'Validation failed.' }
terraform plan -out=create.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Create plan failed.' }
terraform show create.tfplan
```

Inspect the plan. It must create only the intended resource group/server. After approval of that exact plan:

```powershell
terraform apply create.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Creation failed; inspect state and Azure before retrying.' }
$SourceId = terraform output -raw server_id
if ($LASTEXITCODE -ne 0) { throw 'Cannot read source ID.' }
$SourceFqdn = terraform output -raw server_fqdn
if ($LASTEXITCODE -ne 0) { throw 'Cannot read source FQDN.' }
terraform state pull | Set-Content -Encoding utf8 (Join-Path $Evidence 'before-rename.tfstate')
if ($LASTEXITCODE -ne 0) { throw 'Cannot back up source state.' }
```

Applying a saved plan does not prompt again. Treat reviewing and executing it as the approval boundary. Do not use `-auto-approve` for this manual procedure.

### 7.2 Establish baseline data and metadata

1. Confirm Azure reports the source as `Ready`. Record ID, FQDN, location, SKU, storage, backup, HA, identity/CMK, networking, and replication role/source where applicable.
2. Use an already-approved data-plane network path. If connectivity is unavailable, get approval for the smallest necessary network change; do not open the server to all addresses.
3. Use an installed MySQL client with TLS and an interactive password prompt:

```powershell
mysql --host=$SourceFqdn --user=renameadmin --password --ssl-mode=REQUIRED
```

In that client, create test-only data:

```sql
CREATE DATABASE rename_probe;
CREATE TABLE rename_probe.evidence (id INT PRIMARY KEY, note VARCHAR(100));
INSERT INTO rename_probe.evidence VALUES (1, 'before rename');
SELECT * FROM rename_probe.evidence ORDER BY id;
SELECT @@server_uuid, @@version;
```

Save the results privately. The SQL database is deliberately not a separately managed Terraform resource in this standalone test. Data persistence is stronger evidence of in-place behavior than an ARM ID check alone.

### 7.3 Enable rename without changing the name

Edit only `rename_enabled` in `terraform.tfvars` to `true`, keeping `server_name` unchanged.

```powershell
terraform validate
if ($LASTEXITCODE -ne 0) { throw 'Validation failed.' }
terraform plan -out=enable-rename.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Enablement plan failed.' }
terraform show enable-rename.tfplan
```

After verifying it contains only the opt-in change:

```powershell
terraform apply enable-rename.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Cannot record rename opt-in.' }
terraform state show azurerm_mysql_flexible_server.test
if ($LASTEXITCODE -ne 0) { throw 'Cannot inspect source state.' }
```

Confirm `rename_enabled = true` is recorded and the ID/FQDN are unchanged. This step is intentionally separate: persisting opt-in first enables the missing-source safety check even if a later provider process terminates before returning updated state.

### 7.4 Check plan-only negative cases

Use `prevent_destroy = true` and **do not apply these plans**:

| Temporary configuration | Expected result |
| --- | --- |
| `rename_enabled = false` plus a different name | Replacement is required and `prevent_destroy` blocks it |
| Opt-in true, different name, and a tag/SKU/location/resource-group change | Planning error requiring separate changes |
| Empty, too-short, too-long, uppercase, or otherwise invalid name | Name validation error; case-only changes also must not submit a rename |
| Same name and opt-in true | No rename |

Restore the source name, opt-in true, and original properties after each negative case. Run `terraform validate` before each `terraform plan`; validation failures are the expected result for invalid inputs. Unknown opt-in and destination-name values are covered by the automated planning test.

### 7.5 Plan the source-to-destination rename

Edit only `server_name` in `terraform.tfvars` to the value of `$DestinationName`.

```powershell
terraform validate
if ($LASTEXITCODE -ne 0) { throw 'Validation failed.' }
terraform plan -out=rename.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Rename plan failed.' }
terraform show rename.tfplan
$PlanText = terraform show -json rename.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Cannot inspect saved rename plan.' }
$Plan = ($PlanText -join "`n") | ConvertFrom-Json
$Change = @($Plan.resource_changes | Where-Object address -eq 'azurerm_mysql_flexible_server.test')
if ($Change.Count -ne 1 -or ($Change[0].change.actions -join ',') -ne 'update') {
  throw 'Expected an in-place server update, not replacement or creation.'
}
if (-not $Change[0].change.after_unknown.id -or -not $Change[0].change.after_unknown.fqdn) {
  throw 'The destination ID and FQDN must be unknown until apply.'
}
```

Inspect **all** resource changes, not only the server. Stop for unexpected creates, deletes, replacements, or dependency changes. Do not bypass `prevent_destroy`, refresh, or locking to make the plan succeed.

### 7.6 Apply exactly once and verify completion

After approval of the reviewed rename plan:

```powershell
terraform apply rename.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Rename outcome requires investigation; use section 11, not a blind retry.' }
$DestinationId = terraform output -raw server_id
if ($LASTEXITCODE -ne 0) { throw 'Cannot read destination ID.' }
$DestinationFqdn = terraform output -raw server_fqdn
if ($LASTEXITCODE -ne 0) { throw 'Cannot read destination FQDN.' }
terraform state show azurerm_mysql_flexible_server.test
if ($LASTEXITCODE -ne 0) { throw 'Cannot inspect renamed state.' }
terraform state pull | Set-Content -Encoding utf8 (Join-Path $Evidence 'after-rename.tfstate')
if ($LASTEXITCODE -ne 0) { throw 'Cannot back up renamed state.' }
```

Verify all of the following:

1. Apply waited for completion instead of returning success on initial acceptance.
2. The same Terraform address exists; no second server address claims the destination.
3. `id` ends in `$DestinationName`; `name` and FQDN identify the destination.
4. The Terraform resource identity's name matches the destination. The automated acceptance test checks this; the private raw state can also be inspected without editing it.
5. Subscription, resource group, location, administrator, SKU, storage, and other intended settings remain unchanged.
6. The service reports the destination as `Ready`; verify it through the Azure portal or a read-only CLI command. Record old-ID behavior without treating an old-ID 404 alone as proof of success.
7. Connect to `$DestinationFqdn` with the same approved MySQL credentials, read the existing row, insert a second row, and read both:

```sql
SELECT * FROM rename_probe.evidence ORDER BY id;
INSERT INTO rename_probe.evidence VALUES (2, 'after rename');
SELECT * FROM rename_probe.evidence ORDER BY id;
SELECT @@server_uuid, @@version;
```

Any unexplained data/settings change is a failure. Server UUID is useful correlation evidence, not a universal identity guarantee for every failover/topology.

### 7.7 Prove refresh and plan convergence

```powershell
terraform validate
if ($LASTEXITCODE -ne 0) { throw 'Validation failed.' }
terraform plan -detailed-exitcode
$PlanExit = $LASTEXITCODE
if ($PlanExit -ne 0) {
  throw "Expected no changes after rename; exit code was $PlanExit (1=error, 2=changes)."
}
```

Do not accept an endless name/source-ID diff, another rename, an unexpected identity-change error, or replacement.

### 7.8 Repeat rename and perform a normal update

1. Change only `server_name` to `$FinalName`.
2. Repeat sections 7.5-7.7 with a new plan filename and fresh state backup.
3. Verify both SQL rows at the final endpoint.
4. In a **separate apply**, add a harmless test tag to the server resource.
5. Validate, plan, review, and apply that tag change. Confirm it targets the final ID and does not rename again.
6. Run the no-change plan check again.
7. Optionally disable `rename_enabled` in a separate apply without changing `name`; verify the remote server is unchanged. Do not change its name while opt-in is false.

### 7.9 Clean up the disposable resources

Do not clean up while a rename outcome is uncertain. First confirm the final server identity and absence of an in-flight operation.

1. Back up the final state and record the SQL/metadata results.
2. Change the **literal** `prevent_destroy` value to `false` in this test configuration only.
3. Review a destroy plan:

```powershell
terraform validate
if ($LASTEXITCODE -ne 0) { throw 'Validation failed.' }
terraform plan -destroy -out=cleanup.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Cleanup plan failed.' }
terraform show cleanup.tfplan
```

After explicit approval that every listed resource is disposable:

```powershell
terraform apply cleanup.tfplan
if ($LASTEXITCODE -ne 0) { throw 'Cleanup incomplete; retain state and investigate.' }
terraform state list
```

Verify the final server and intended test resource group are gone through a read-only Azure check. A Go test reporting failure does not prove cleanup succeeded. Retain evidence/state according to policy, then remove only known test artifacts; do not recursively delete a shared workspace.

## 8. Run the live Go acceptance test

This is an **additional** live run: it creates its own resources and is not the manual configuration from section 7.

### 8.1 Configure an isolated acceptance shell

Use a fresh shell with the correct Go and Terraform binaries. Do not reuse the manual development override or its `TF_VAR_*` settings.

```powershell
Set-Location 'C:\Work\repos\terraform-provider-azurerm'
Remove-Item Env:\TF_CLI_CONFIG_FILE -ErrorAction SilentlyContinue
Remove-Item Env:\TF_REATTACH_PROVIDERS -ErrorAction SilentlyContinue
Remove-Item Env:\GO_AZURE_SDK_SKIP_POLLING_DELAY -ErrorAction SilentlyContinue

$env:TF_ACC = '1'
$env:ARM_MYSQL_RENAME_ENABLED = 'true'
$env:TF_SCHEMA_PANIC_ON_ERROR = '1'
$env:TF_ACC_TERRAFORM_PATH = (Get-Command terraform).Source
$env:ARM_CLIENT_ID = '<approved-existing-test-client-id>'
$env:ARM_CLIENT_SECRET = Read-Host 'Approved existing test client secret' -MaskInput
$env:ARM_SUBSCRIPTION_ID = '<approved-subscription-id>'
$env:ARM_TENANT_ID = '<approved-tenant-id>'
$env:ARM_USE_CLI = 'false'
$env:ARM_TEST_LOCATION = '<approved-primary-test-region>'
$env:ARM_TEST_LOCATION_ALT = '<approved-secondary-test-region>'
$env:ARM_TEST_LOCATION_ALT2 = '<approved-rename-enabled-region>'
```

Use your approved alternative authentication method if required instead of adding a client secret. Never print or publish the credential variables.

**Important:** The existing MySQL test template uses `data.Locations.Ternary`, which is **`ARM_TEST_LOCATION_ALT2`**. Setting only `ARM_TEST_LOCATION` to a rename-enabled region is insufficient.

The explicit acceptance Terraform path prevents the harness from selecting a different installed/downloaded CLI. Use Terraform 1.12 or newer for the rename test's identity assertions.

### 8.2 Run only the new live test

```powershell
go test -p 6 .\internal\services\mysql `
  -run '^TestAccMySqlFlexibleServer_rename$' -count=1 -v -timeout 180m `
  '-ldflags=-X=github.com/hashicorp/terraform-provider-azurerm/version.ProviderVersion=acc'
if ($LASTEXITCODE -ne 0) {
  throw 'Live acceptance failed. Inspect the operation, resources, and cleanup before rerunning.'
}
```

Equivalent repository Make entry point, in the supported POSIX shell after setting the same live environment:

```bash
make acctests SERVICE='mysql' TESTARGS='-run=^TestAccMySqlFlexibleServer_rename$ -count=1' TESTTIMEOUT='180m'
```

Expected sequence: create source; rename to `-renamed`; no-change plan; import verification; rename to `-again`; import verification; cleanup. The shared test harness also rejects replacement of the tested resource and refreshes after apply.

Expected result is an actual **PASS**, not SKIP. The test ignores password and provider-only opt-in differences during import verification because neither is returned by Azure.

The test does not prove SQL data persistence, real failure recovery, HA/CMK/private networking, or replica propagation. Those require the manual matrix.

### 8.3 Run the agreed existing regression tests

After separate approval of their cost, runtime, and regions, run the relevant existing cases:

```powershell
go test -p 6 .\internal\services\mysql `
  -run '^TestAccMySqlFlexibleServer_(basic|completeUpdate|updateMaintenanceWindow|updateSku|updateHA|replica|geoRestore)$|^TestAccMysqlFlexibleServer_resourceIdentity$' `
  -count=1 -v -timeout 360m `
  '-ldflags=-X=github.com/hashicorp/terraform-provider-azurerm/version.ProviderVersion=acc'
if ($LASTEXITCODE -ne 0) { throw 'Existing live regression tests failed.' }
```

The geo-restore test has an additional `ARM_GEO_RESTORE_LOCATION` prerequisite, and other tests may impose environment-specific requirements. Read each selected test before running it; record skips honestly. Expand to private-networking, CMK, child-resource, and wider MySQL tests as maintainers request. Do not run all services' acceptance tests.

When finished, remove live opt-in flags and secrets from this shell, or close it:

```powershell
Remove-Item Env:\TF_ACC, Env:\ARM_MYSQL_RENAME_ENABLED, Env:\ARM_CLIENT_SECRET `
  -ErrorAction SilentlyContinue
```

## 9. Test import and an existing server

Do not create a second resource to "take over" an already-managed server.

1. For an existing Terraform-managed server, back up its state, keep its resource address, load the development provider, and reconcile existing drift **before** enabling rename.
2. For an unmanaged disposable server, use an isolated configuration matching its current name and properties. Reference the existing resource group rather than configuring a new group for creation.
3. Confirm no other state/address manages it, then import:

   ```powershell
   $ExistingId = '<confirmed-existing-flexible-server-resource-id>'
   terraform import azurerm_mysql_flexible_server.test $ExistingId
   if ($LASTEXITCODE -ne 0) { throw 'Import failed.' }
   ```

4. Run `terraform validate` and `terraform plan`; resolve configuration/state differences without applying a replacement. `create_mode` is not returned by the API, and importing a replica/restore needs particular care. Consult the resource documentation.
5. Record opt-in with a separate apply, then run the same rename and convergence checks as section 7.
6. Use the acceptance test's import checks to validate the post-rename importer without manually forgetting a healthy managed server just for a test.

## 10. Dependencies and extended topology matrix

Automatic state migration for other Terraform resources is **not implemented**. A changed parent ID/name can plan destructive replacement of databases, replicas, firewall rules, or other references even though the server itself is renamed in place.

### 10.1 Safely separate dependency migration

1. Inventory every managed and external reference: databases, configurations, firewall rules, administrators, identities, CMK access, private endpoints/DNS, role assignments, connection strings, replicas, monitoring, and outputs consumed by other states.
2. Back up all affected states and obtain a reviewed migration plan.
3. Before the server rename, temporarily stop managing affected dependent resources **without destroying them**, using reviewed `removed` blocks with `lifecycle { destroy = false }` or explicitly approved `terraform state rm` operations.
4. Remove or replace configuration references that would otherwise re-create them in the same plan. A `removed` block and a resource block cannot both claim the same address.
5. Plan and apply that state-only separation before the server rename. Check that there are no Azure delete operations.
6. Rename only the intended server. Do not use `-target` as a substitute for a correct dependency migration.
7. Verify each dependent Azure resource and replication link. Determine actual destination IDs; do not guess that every dependent resource moved.
8. Restore matching configurations, update references, and import existing resources at their confirmed IDs. Do not use create to adopt a resource that already exists.
9. Run a full normal plan across affected states. Investigate every proposed destroy/replacement before applying anything.

Keep this procedure manual and explicit. Do not modify another resource's state from inside the parent resource implementation.

### 10.2 Additional live scenarios required for broader claims

Run these independently on disposable resources after confirming service support. They are **not** implemented as additional rename acceptance tests today.

| Scenario | Procedure and required evidence |
| --- | --- |
| Standalone | Full section 7, persistent SQL rows, A-to-B-to-C, import, update, cleanup |
| SameZone/ZoneRedundant HA | Create supported HA topology first; snapshot HA/zones; rename only; verify HA mode, endpoint connectivity, and data |
| Private networking | Create supported private topology first; record subnet/DNS/endpoint links; rename; verify resolution and connectivity from the approved network |
| Customer-managed key | Create CMK/identity configuration first; rename without changing it; verify identity IDs, key references/access, encryption settings, and data |
| Same-region replica rename | Record primary ID, replica source ID and replica-list membership; rename replica; verify membership and primary writes reaching the new replica endpoint |
| Cross-region replica rename | Repeat the replica test with independently enabled, approved regions; verify propagation and replication in both regions |
| Primary with replicas | Separate dependent Terraform management first; rename primary; verify every replica's source ID, membership, metadata, and actual new writes |
| Ordinary CRUD/import after rename | Import final ID, normal update, stable plan, and deletion of the confirmed destination |
| Existing non-rename users | No opt-in: unchanged create/restore/replica/update/import/delete behavior; name replacement only tested with a plan or disposable resources |

Record unsupported scenarios as unsupported or pending, not passed. Backend/operator-level metadata checks must use the service owner's approved runbook and permissions. Do not publish internal control-store rows, hostnames, replication-set identifiers, or operational logs.

## 11. Failure, interruption, and recovery

### 11.1 Exercise failure cases safely

Use offline tests first. Live fault injection needs separate approval and a recovery owner.

| Case | Safe test and expected behavior |
| --- | --- |
| Destination already exists | Use another approved disposable server name; apply must fail without overwriting/adopting it |
| Name unavailable elsewhere | Confirm unavailability; rename must fail and retain source state |
| Feature/API unavailable | In a separately approved environment, capture the real service rejection; do not disable shared feature switches |
| Source not ready | Prefer offline coverage; do not restart/fail over a shared server to manufacture the condition |
| Initial authorization error | Use an approved pre-existing restricted identity; do not revoke shared permissions |
| Polling failure, timeout, or Ctrl-C | Disposable source only; record opt-in first, capture operation information, and inspect its eventual outcome |
| Failed/canceled LRO | Arrange with service owners or retain offline evidence; do not call an undocumented cancellation endpoint |
| Missing operation URL or unexpected destination | Offline coverage today; treat an unverified response as failure, not success |
| Source 404 after out-of-band rename | Provider should retain state and return an explicit error when opt-in is recorded |

The provider uses the SDK's retry behavior. Do not claim exactly one HTTP attempt under every network/server error; the offline checks prove no application-level resubmission for the tested cases, not all transport behaviors.

### 11.2 Diagnose an uncertain result

1. Stop new applies, renames, and cleanup. Preserve the current private state and error output.
2. Record source/destination IDs, UTC time, operation URL if present, and request/correlation IDs.
3. Verify the operation URL belongs to the expected Azure management cloud before using an authenticated client to poll it.
4. Poll the **existing** operation read-only:

   ```powershell
   $OperationUrl = '<verified-operation-url-from-the-error-or-approved-log>'
   az rest --method get --url $OperationUrl
   if ($LASTEXITCODE -ne 0) { throw 'Cannot determine operation status; do not resubmit.' }
   ```

5. Distinguish `Succeeded`, `Failed`, `Canceled`, and still-running states. Read both server IDs and their metadata. A destination that exists might be unrelated; neither name resolving can be transient.
6. If the service rolled back, confirm the original server is `Ready`, no operation remains active, and its data/replication are healthy before considering a new plan.
7. If both names resolve, neither resolves, rollback is incomplete, or ownership is uncertain, involve the service owner. Do not delete either resource or recreate an empty server.

Provider logs can expose sensitive material. Enable detailed logging only when necessary, scope it to the approved run, and store it privately; never attach raw logs or state to a public issue.

### 11.3 Reconcile state after verified late success

Only use this procedure after proving the destination is the original server and the rename finished. Acquire exclusive control of the state/backend and back it up. Do not let another process plan/apply between the commands.

1. Set the existing resource configuration to the **confirmed destination name**, preserving other settings and the same Terraform address.
2. Run `terraform validate`.
3. Forget only the stale source entry and immediately import the confirmed destination:

   ```powershell
   terraform state pull | Set-Content -Encoding utf8 (Join-Path $Evidence 'before-recovery.tfstate')
   if ($LASTEXITCODE -ne 0) { throw 'State backup failed; do not remove anything.' }
   terraform state rm -dry-run azurerm_mysql_flexible_server.test
   if ($LASTEXITCODE -ne 0) { throw 'State address check failed.' }
   ```

   Review the single address and obtain approval for the state change, then:

   ```powershell
   $ConfirmedDestinationId = '<verified-original-server-at-its-new-resource-id>'
   terraform state rm azurerm_mysql_flexible_server.test
   if ($LASTEXITCODE -ne 0) { throw 'State removal failed.' }
   terraform import azurerm_mysql_flexible_server.test $ConfirmedDestinationId
   if ($LASTEXITCODE -ne 0) {
     throw 'Import failed. Do not run apply; repair import/state under the recovery plan.'
   }
   terraform validate
   if ($LASTEXITCODE -ne 0) { throw 'Validation failed.' }
   terraform plan
   if ($LASTEXITCODE -ne 0) { throw 'Recovery plan failed.' }
   ```

4. Inspect ID, identity, FQDN, opt-in, data, and all dependencies. Allow only the understood provider-only opt-in difference if import did not retain it; record that in a separate apply.
5. Require a stable no-change plan before resuming normal management.

State removal is not an Azure delete, but it relinquishes management and is dangerous without immediate verified import. Never use it to initiate a rename or to bypass uncertain ownership.

### 11.4 Rollback limitations

Reverting provider code does not undo an Azure rename. Renaming back is a new approved service operation, not an automatic rollback guarantee. An older provider may reject the new `rename_enabled` argument or handle newer state differently. Prefer a maintainer-supported forward-fix release; do not downgrade production state blindly.

If an out-of-band deletion was intentional, verify deletion and explicitly reconcile its stale state rather than repeatedly clearing the rename safety error.

## 12. Prepare and open the feature PR

These are **future user-authorized Git/GitHub actions**, not actions performed by creating this document.

### 12.1 Obtain the required design and API approval

Before submitting, ask maintainers to review:

1. The `2025-12-30` non-preview API choice and explicit opt-in/default replacement behavior.
2. The use of `2023-12-30` generated server models with a separately configured `2025-12-30` rename client; a generated SDK update or split API-update PR may be required.
3. Mutable resource identity, unknown ID/FQDN planning with the narrow protocol adapter, plan-time destination-name requirement, and missing-source refresh behavior.
4. Separate-property-change restrictions and the documented manual dependency/recovery procedure.
5. Live support evidence for the `2025-12-30` endpoint and the scope of topology claims.

The rename client uses a non-preview API version. Keep the repository's [preview API version check](../../.github/workflows/pr-check-preview-api-version.yaml) intact; a green or skipped check does not validate live rename support or the compatibility of the reused generated models.

### 12.2 Prepare a feature branch without losing local changes

Once ready and explicitly authorized to create a branch:

```powershell
Set-Location $Repo
git status --short
git switch -c mysql-flexible-server-rename
if ($LASTEXITCODE -ne 0) { throw 'Branch creation failed; inspect existing branches instead of forcing.' }
```

Do not submit from `main`. If that branch name already exists, inspect it and choose an appropriate unused name; do not reset it.

Run the focused tests and full gates on the feature revision. If upstream has moved, integrate changes using the project's reviewed workflow, preserving local work; do not force-push or rebase shared history by default.

### 12.3 Prepare the changelog and commit only intended files

In the supported POSIX shell:

```bash
make changelog-types
make changelog TYPE=resource-enhancement BODY='`azurerm_mysql_flexible_server` - add opt-in support for renaming an existing server'
```

Without a PR number this creates a placeholder entry, normally `.changelog/0.hcl`. Check that the path is not already owned by another change; resolve any collision rather than overwriting it. Once the PR exists, rename the entry to its actual number and update any entry metadata required by the tooling.

In PowerShell, review/stage the explicit feature list, the intended runbook/index changes, and the actual changelog file:

```powershell
git add -- internal\services\mysql\client\client.go `
  internal\services\mysql\mysql_flexible_server_resource.go `
  internal\services\mysql\mysql_flexible_server_rename.go `
  internal\services\mysql\mysql_flexible_server_rename_test.go `
  internal\services\mysql\mysql_flexible_server_rename_resource_test.go `
  website\docs\r\mysql_flexible_server.html.markdown `
  contributing\topics\mysql-flexible-server-rename-testing.md
git diff --cached --check
git diff --cached --stat
git diff --cached
```

Stage the contributor index only if its intended link is included, and stage the inspected changelog separately. Do not use `git add .`, include local binaries/state/evidence, or publish internal CLI/backend notes.

After explicit commit approval, create a commit with a meaningful message and the required contribution/AI trailers. For the implementation associated with this runbook:

```powershell
git commit `
  -m 'azurerm_mysql_flexible_server: add opt-in in-place rename' `
  -m "Co-authored-by: Copilot App <223556219+Copilot@users.noreply.github.com>`nCopilot-Session: 798b7031-ede9-4ebe-9d08-4c99d5d5a0f9"
if ($LASTEXITCODE -ne 0) { throw 'Commit failed; inspect the failure before retrying.' }
```

Then rerun checks if staging or generation changed code. Subsequent independently authored work should use its own applicable attribution rather than copying this session trailer.

### 12.4 Push to the intended fork

1. Verify `git remote -v`; the inspected checkout's `origin` points to `hashicorp/terraform-provider-azurerm`.
2. Use your existing authorized fork. If none exists, create one through the approved GitHub workflow.
3. Add a separate fork remote without rewriting the existing upstream remote:

```powershell
git remote add fork https://github.com/<your-github-user>/terraform-provider-azurerm.git
git remote -v
```

If `fork` already exists, verify its URL instead of replacing it. After explicit push approval:

```powershell
git push -u fork mysql-flexible-server-rename
if ($LASTEXITCODE -ne 0) { throw 'Push failed; do not force-push.' }
```

### 12.5 Open the PR using the current repository template

Use GitHub's compare/open-PR UI or the app's Create Pull Request action. Select upstream `hashicorp/terraform-provider-azurerm`, base `main`, and the authorized feature branch in your fork. Do not trigger release workflows.

Use a title such as **`azurerm_mysql_flexible_server - add opt-in in-place rename support`**.

Start from the actual [PR template](../../.github/pull_request_template.md), not a generic replacement. Preserve its headings, comments, checklist items, and order. In particular:

| Template area | Required content |
| --- | --- |
| Community Note | Leave unchanged |
| Description | Lifecycle behavior, API version/client choice, compatibility opt-in, state and dependency limits |
| PR Checklist / Changes to existing Resource | Check only completed items |
| Testing | Exact commands **and results**, revision, environment class, live passes/skips, and sanitized evidence |
| Change Log | Actual `.changelog/<PR-number>.hcl`; do not claim this user-facing feature needs no changelog |
| PR type | Appropriate feature/enhancement and any agreed compatibility classification |
| Related Issue(s) | Real issue only; do not leave the `Fixes #0000` placeholder or invent one |
| AI Assistance Disclosure | Disclose AI-assisted implementation, tests, documentation, and any AI-assisted review replies |
| Rollback Plan | Forward-fix/provider-revert process plus separate Azure/state recovery; no promise of automatic rename reversal |
| Changes to Security Controls | Explain changed identity/reference behavior and any actual security-control impact |

The current PR-opening guide asks for review-ready work rather than draft/incomplete feature PRs. Discuss unresolved API-client/lifecycle decisions before submission.

After opening, obtain the real PR number, rename the changelog entry, and push that narrowly scoped follow-up. Include only sanitized summaries of private live evidence.

### 12.6 Complete CI and review

```powershell
$PrNumber = '<actual-pr-number>'
gh pr checks $PrNumber --repo hashicorp/terraform-provider-azurerm
gh pr view $PrNumber --repo hashicorp/terraform-provider-azurerm --web
```

1. Resolve failures in applicable combined PR checks, generation/schema/formatting, lint, documentation, links, and changelog checks.
2. Confirm API-version and client-model compatibility approval explicitly; skipped checks are not approvals.
3. Request the maintainer-run acceptance tests through the documented project process. Fork CI generally cannot use protected Azure credentials.
4. Provide the opt-in and `ARM_TEST_LOCATION_ALT2` requirements; otherwise the new live test can skip or run in the wrong region.
5. Address every review comment, rerun affected tests, and update the PR's evidence when the implementation changes.
6. Wait for required approval and checks. Maintainers use squash-and-merge per the [merge guide](maintainer-merging.md).

Do not claim production readiness from offline tests, CI skips, or merge alone.

## 13. Backend rollout and enablement

Skip actual deployment if the approved test environment already runs the necessary service code and feature configuration.

If service work is required:

1. Have the service owner identify the approved backend revision implementing rename for the provider's API version and topology.
2. Confirm API contract publication/routing, source/destination validation, asynchronous operation behavior, rollback, and replication propagation prerequisites.
3. Use the service team's approved build pipeline and validate its immutable artifact.
4. Obtain a scoped rollout/change approval, health criteria, rollback owner, and test window.
5. Deploy to the approved test/canary scope using that team's runbook. Canary is still a production environment; it is not automatically safe.
6. Apply only approved feature enablement. Do not copy global override commands from historical notes or infer subscription scope from a command's description.
7. Verify the actual deployed version, configuration scope, and both regions for cross-region replication tests.
8. Run the native Terraform tests, data-plane checks, and authorized service metadata checks from this document.
9. Remove temporary overrides through the approved procedure and retain private audit evidence.

Backend pipeline names, rollout commands, and feature-switch values are intentionally not reproduced here: this public provider checkout does not define or authorize them. Use the current service-owned canary runbook, not stale example outputs. This document does not authorize changes to any Azure subscription or service control plane.

## 14. Provider packaging, release, and consumer rollout

### 14.1 Maintainer-only release preparation

After merge and acceptance:

1. An authorized release owner selects the release version, branch, and exact commit containing the feature.
2. Prepare release notes/schema/version updates using the current repository release process. `make release-changelog` changes release artifacts; it is not a contributor test command.
3. Run the release build and inspect artifacts, checksums, signing/provenance, metadata, and supported platform coverage.
4. Confirm the selected artifact includes this change and the resource schema exposes `rename_enabled`.

The inspected [build workflow](../../.github/workflows/build.yml) has `workflow_dispatch` and push triggers limited to `main`/`release/**` changes to `version/VERSION`. Do not assume every feature-branch push produces a distributable release artifact.

### 14.2 Staging and production promotion

Only an authorized release operator performs these actions in the repository's **Release Provider** workflow:

1. Select the tested release version, release branch, and exact commit SHA.
2. Use the `staging` release channel first.
3. Verify staging artifacts through the approved installation route and run a smoke test with the package rather than the local development binary.
4. Confirm acceptance, service readiness, documentation, and rollback/forward-fix readiness.
5. Request `production` promotion with the same approved inputs.
6. Obtain the protected deployment approval and verify final promotion completion.
7. Confirm the published release and Registry package/documentation actually contain the feature.

The inspected [release workflow](../../.github/workflows/provider-release.yml) dispatches Common Release Tooling promotion. Its inputs are `release-version`, `release-branch`, `git-commit-sha`, and `release-channel`. Production is approval-gated. Do not manually invent tags, upload unsigned binaries, or invoke this workflow from a contributor fork as a substitute.

### 14.3 Test normal customer installation

Use a separate approved smoke-test workspace and the actual published version; do not invent a release number.

1. Remove the development override **from this shell**:

   ```powershell
   Remove-Item Env:\TF_CLI_CONFIG_FILE -ErrorAction SilentlyContinue
   Remove-Item Env:\TF_REATTACH_PROVIDERS -ErrorAction SilentlyContinue
   ```

2. Verify no global CLI configuration, plugin mirror, or CI setting still redirects AzureRM to a development build.
3. Set `required_providers.azurerm.version` to an exact tested version that includes the feature.
4. Initialize the normal installation and inspect the selected package:

   ```powershell
   terraform init -upgrade
   if ($LASTEXITCODE -ne 0) { throw 'Published provider installation failed.' }
   terraform providers
   terraform providers schema -json
   ```

5. Confirm `.terraform.lock.hcl` contains the intended provider version and normal package hashes. Commit the lock file in the consumer project according to its policy, never the provider test state.
6. Repeat schema validation, standalone rename, no-change plan, and cleanup using the released package.
7. Roll out the provider upgrade to non-production consumers first. Upgrade alone must not rename servers; rename requires an explicit opt-in and reviewed name change.
8. For each intended production rename, repeat the service-readiness, state/dependency, endpoint, approval, and recovery gates.
9. Expand gradually, monitor results, and stop on identity drift, data/replication failures, or ambiguous completion.

Do not leave development overrides enabled in customer/production automation. Removing an override is not sufficient if another configuration or mirror still supplies the custom binary.

## 15. Evidence checklist and completion criteria

Keep private detailed evidence and publish only sanitized summaries.

| Evidence | Record |
| --- | --- |
| Source | Base and feature commit, dirty-state status, exact diff, PR link |
| Tools | Go/Terraform versions, OS/architecture, provider build hash or released package/lock entry |
| Offline validation | Exact commands, UTC timestamp, pass/fail, and any skipped live test |
| Full repository gates | Make/CI results, generation differences, API-client compatibility approval |
| Live scope | Approved environment class, region/topology, permission/change approval reference |
| Rename | Before/after IDs and FQDNs privately, operation result, state identity, update-only plan |
| Data and replication | Persistent rows, post-rename writes, replica propagation and required settings |
| Recovery | Which faults were exercised, actual final service outcome, state reconciliation, no blind retries |
| Cleanup | Confirmed final-ID deletion and any remaining resource owner |
| Release | Merge commit, build/provenance, staging/production approvals, published version, package smoke test |

Completion gates:

- [ ] Offline rename and MySQL validation/parser suites pass on the final revision.
- [ ] Full repository/CI gates pass or an authorized maintainer explicitly resolves a documented exception.
- [ ] Preview API and lifecycle design are approved.
- [ ] A real Terraform CLI run uses the intended development binary and performs no server replacement.
- [ ] Live standalone acceptance, identity, import, convergence, normal update, and cleanup pass without skips.
- [ ] SQL data persistence and endpoint behavior are verified.
- [ ] Every topology claimed as supported has its own live evidence.
- [ ] Failure/late-success recovery and dependency migration have an approved operational procedure.
- [ ] PR template, changelog, AI disclosure, and sanitized evidence are complete.
- [ ] Maintainer merge and authorized release promotion complete.
- [ ] The published package, without development overrides, passes the smoke test.
- [ ] Production consumers receive the provider and backend capability through separately approved rollouts.

Until those live/release gates are met, describe this as **locally implemented and offline-tested**, not generally available or production-certified.

## References

- [MySQL Flexible Server resource documentation](../../website/docs/r/mysql_flexible_server.html.markdown)
- [Building and locally loading the provider](building-the-provider.md)
- [Running tests](running-the-tests.md)
- [Opening a PR](guide-opening-a-pr.md)
- [Resource identity](guide-resource-identity.md)
- [Maintainer merging and changelogs](maintainer-merging.md)
- [Development overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)
- [Terraform state commands](https://developer.hashicorp.com/terraform/cli/commands/state)
- [Terraform import](https://developer.hashicorp.com/terraform/cli/commands/import)
- [Removing resources without destroying them](https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources)
