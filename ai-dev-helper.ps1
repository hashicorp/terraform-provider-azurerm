#Requires -Version 5.1

<#
.SYNOPSIS
    Simple AI Agent Helper for Terraform AzureRM Provider Development

.DESCRIPTION
    A streamlined, single-file solution for local AI-enhanced development.
    No dependencies, no modules, just simple functions for common tasks.

.EXAMPLE
    # Import the functions
    . .\ai-dev-helper.ps1
    
    # Run a quick test
    Test-AcceptanceTest -Service "cdn" -Test "basic"
    
    # Build provider
    Build-Provider
    
    # Quick validation
    Validate-Environment
#>

# Global configuration
$Script:TerraformProviderPath = $PSScriptRoot
$Script:VerboseOutput = $false

#region Core Functions

function Write-StatusMessage {
    param([string]$Message, [string]$Type = "Info")
    
    $timestamp = Get-Date -Format "HH:mm:ss"
    switch ($Type) {
        "Success" { Write-Host "[$timestamp] ✓ $Message" -ForegroundColor Green }
        "Warning" { Write-Host "[$timestamp] ⚠ $Message" -ForegroundColor Yellow }
        "Error"   { Write-Host "[$timestamp] ✗ $Message" -ForegroundColor Red }
        default   { Write-Host "[$timestamp] ℹ $Message" -ForegroundColor Cyan }
    }
}

function Test-Prerequisites {
    Write-StatusMessage "Checking prerequisites..."
    
    # Check Go
    try {
        $goVersion = go version 2>$null
        if ($goVersion -match "go(\d+\.\d+)") {
            Write-StatusMessage "Go version: $($matches[1])" "Success"
        }
    } catch {
        Write-StatusMessage "Go not found - install Go 1.22+" "Error"
        return $false
    }
    
    # Check if we're in the right directory
    if (-not (Test-Path "go.mod")) {
        Write-StatusMessage "Not in Terraform provider directory" "Error"
        return $false
    }
    
    # Check Azure credentials for testing
    $azureVars = @("ARM_SUBSCRIPTION_ID", "ARM_CLIENT_ID", "ARM_CLIENT_SECRET", "ARM_TENANT_ID")
    $missingVars = $azureVars | Where-Object { -not $env:$_ }
    
    if ($missingVars) {
        Write-StatusMessage "Missing Azure credentials: $($missingVars -join ', ')" "Warning"
    } else {
        Write-StatusMessage "Azure credentials configured" "Success"
    }
    
    return $true
}

#endregion

#region Build Functions

function Build-Provider {
    <#
    .SYNOPSIS
        Build the Terraform AzureRM Provider
    #>
    
    Write-StatusMessage "Building provider..."
    
    try {
        $result = go build -o terraform-provider-azurerm.exe . 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-StatusMessage "Provider built successfully" "Success"
            return $true
        } else {
            Write-StatusMessage "Build failed: $result" "Error"
            return $false
        }
    } catch {
        Write-StatusMessage "Build error: $_" "Error"
        return $false
    }
}

function Test-Build {
    <#
    .SYNOPSIS
        Quick build test (no output file)
    #>
    
    Write-StatusMessage "Testing build..."
    
    try {
        $result = go build . 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-StatusMessage "Build test passed" "Success"
            return $true
        } else {
            Write-StatusMessage "Build test failed: $result" "Error"
            return $false
        }
    } catch {
        Write-StatusMessage "Build test error: $_" "Error"
        return $false
    }
}

#endregion

#region Testing Functions

function Test-AcceptanceTest {
    <#
    .SYNOPSIS
        Run acceptance tests for a specific service and test
        
    .PARAMETER Service
        Azure service name (e.g., 'cdn', 'compute', 'storage')
        
    .PARAMETER Test
        Test name pattern (e.g., 'basic', 'requiresImport', 'update')
        
    .PARAMETER Resource
        Specific resource name (optional)
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Service,
        
        [Parameter(Mandatory)]
        [string]$Test,
        
        [string]$Resource = ""
    )
    
    if (-not (Test-Prerequisites)) {
        return $false
    }
    
    # Build test path
    $testPath = "./internal/services/$Service"
    if (-not (Test-Path $testPath)) {
        Write-StatusMessage "Service path not found: $testPath" "Error"
        return $false
    }
    
    # Build test pattern
    $testPattern = if ($Resource) {
        "TestAcc${Resource}_${Test}"
    } else {
        "TestAcc*_${Test}"
    }
    
    Write-StatusMessage "Running tests: $testPattern in $Service"
    
    try {
        $env:TF_ACC = "1"
        $cmd = "go test -v $testPath -run `"$testPattern`" -timeout 60m"
        Write-StatusMessage "Command: $cmd"
        
        Invoke-Expression $cmd
        
        if ($LASTEXITCODE -eq 0) {
            Write-StatusMessage "Tests passed" "Success"
            return $true
        } else {
            Write-StatusMessage "Tests failed" "Error"
            return $false
        }
    } catch {
        Write-StatusMessage "Test execution error: $_" "Error"
        return $false
    }
}

function Test-UnitTests {
    <#
    .SYNOPSIS
        Run unit tests for a service
    #>
    param([string]$Service)
    
    Write-StatusMessage "Running unit tests for $Service..."
    
    try {
        $testPath = "./internal/services/$Service"
        $result = go test $testPath 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-StatusMessage "Unit tests passed" "Success"
            return $true
        } else {
            Write-StatusMessage "Unit tests failed: $result" "Error"
            return $false
        }
    } catch {
        Write-StatusMessage "Unit test error: $_" "Error"
        return $false
    }
}

#endregion

#region Development Helpers

function Find-Resource {
    <#
    .SYNOPSIS
        Find resource files by name pattern
    #>
    param([string]$Pattern)
    
    Write-StatusMessage "Searching for resources matching: $Pattern"
    
    $results = Get-ChildItem -Path "internal/services" -Recurse -Filter "*$Pattern*" | 
        Where-Object { $_.Name -match "resource|data_source" -and $_.Extension -eq ".go" }
    
    if ($results) {
        $results | ForEach-Object {
            Write-Host "  $($_.FullName)" -ForegroundColor Gray
        }
        return $results
    } else {
        Write-StatusMessage "No resources found matching: $Pattern" "Warning"
        return $null
    }
}

function Show-ServiceStructure {
    <#
    .SYNOPSIS
        Show the structure of a service directory
    #>
    param([string]$Service)
    
    $servicePath = "internal/services/$Service"
    if (-not (Test-Path $servicePath)) {
        Write-StatusMessage "Service not found: $Service" "Error"
        return
    }
    
    Write-StatusMessage "Structure for service: $Service"
    Get-ChildItem $servicePath | ForEach-Object {
        $type = if ($_.PSIsContainer) { "DIR " } else { "FILE" }
        Write-Host "  $type $($_.Name)" -ForegroundColor Gray
    }
}

function Validate-Environment {
    <#
    .SYNOPSIS
        Quick environment validation
    #>
    
    Write-StatusMessage "Environment Validation Report"
    Write-Host "=" * 40
    
    # Prerequisites
    $prereqOk = Test-Prerequisites
    
    # Build test
    $buildOk = Test-Build
    
    # Check for common issues
    $issuesFound = @()
    
    if (-not (Test-Path ".git")) {
        $issuesFound += "Not in a Git repository"
    }
    
    if (-not (Test-Path "internal/services")) {
        $issuesFound += "Services directory not found"
    }
    
    # Summary
    Write-Host ""
    if ($prereqOk -and $buildOk -and $issuesFound.Count -eq 0) {
        Write-StatusMessage "Environment is ready for development!" "Success"
    } else {
        Write-StatusMessage "Issues found:" "Warning"
        $issuesFound | ForEach-Object { Write-Host "  - $_" -ForegroundColor Yellow }
    }
}

#endregion

#region AI Agent Helpers

function Generate-ResourceTest {
    <#
    .SYNOPSIS
        Generate basic test structure for a new resource
    #>
    param(
        [Parameter(Mandatory)]
        [string]$Service,
        
        [Parameter(Mandatory)]
        [string]$ResourceName
    )
    
    $testTemplate = @"
func TestAccAzureRM${ResourceName}_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_${Service}_${ResourceName}", "test")
    r := ${ResourceName}Resource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

func TestAccAzureRM${ResourceName}_requiresImport(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_${Service}_${ResourceName}", "test")
    r := ${ResourceName}Resource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.RequiresImportErrorStep(r.requiresImport),
    })
}

func (r ${ResourceName}Resource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(\`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_${Service}_${ResourceName}" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
\`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
"@

    Write-StatusMessage "Generated test template for ${Service}/${ResourceName}:"
    Write-Host $testTemplate -ForegroundColor Gray
}

function Quick-Commands {
    <#
    .SYNOPSIS
        Show common commands for quick reference
    #>
    
    Write-Host @"

🚀 AI Development Helper - Quick Commands

BUILD & TEST:
  Build-Provider                              # Build the provider
  Test-Build                                  # Quick build test
  Test-AcceptanceTest -Service cdn -Test basic  # Run acceptance tests
  Test-UnitTests -Service cdn                 # Run unit tests

DEVELOPMENT:
  Find-Resource -Pattern "frontdoor"         # Find resources
  Show-ServiceStructure -Service cdn         # Show service structure
  Generate-ResourceTest -Service cdn -ResourceName profile  # Generate test template

VALIDATION:
  Validate-Environment                       # Check environment
  Test-Prerequisites                         # Check prerequisites

EXAMPLES:
  Test-AcceptanceTest -Service cdn -Test basic -Resource "CdnFrontDoorProfile"
  Find-Resource -Pattern "cdn"
  Show-ServiceStructure -Service compute

"@ -ForegroundColor Cyan
}

#endregion

# Initialize
Write-StatusMessage "AI Development Helper loaded! Type 'Quick-Commands' for help." "Success"

# Auto-validate if in correct directory
if (Test-Path "go.mod") {
    $goMod = Get-Content "go.mod" | Select-Object -First 1
    if ($goMod -match "terraform-provider-azurerm") {
        Write-StatusMessage "Terraform AzureRM Provider detected" "Success"
    }
}
