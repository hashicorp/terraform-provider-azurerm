# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# follows https://learn.microsoft.com/en-us/azure-stack/hci/deploy/deployment-arc-register-server-permissions?tabs=powershell
# the script is from https://github.com/Azure/Edge-infrastructure-quickstart-template/blob/c42d671bd8464071c3b8dfc9d599bef1631b22a2/modules/hci-provisioners/hci-server/connect.ps1
param(
    $userName,
    $password,
    $authType,
    $ip, $port,
    $subscriptionId, $resourceGroupName, $region, $tenant, $servicePrincipalId, $servicePrincipalSecret, $expandC
)

$script:ErrorActionPreference = 'Stop'
echo "Start to connect Arc server!"
$count = 0

if ($authType -eq "CredSSP") {
    try {
        Enable-WSManCredSSP -Role Client -DelegateComputer $ip -Force
    }
    catch {
        echo "Enable-WSManCredSSP failed"
    }
}
for ($count = 0; $count -lt 6; $count++) {
    try {
        $secpasswd = ConvertTo-SecureString $password -AsPlainText -Force
        $cred = New-Object System.Management.Automation.PSCredential -ArgumentList ".\$username", $secpasswd
        $session = New-PSSession -ComputerName $ip -Port $port -Authentication $authType -Credential $cred

        Invoke-Command -Session $session -ScriptBlock {
            Param ($subscriptionId, $resourceGroupName, $region, $tenant, $servicePrincipalId, $servicePrincipalSecret)
            $script:ErrorActionPreference = 'Stop'

            function Install-ModuleIfMissing {
                param(
                    [Parameter(Mandatory = $true)]
                    [string]$Name,
                    [string]$Repository = 'PSGallery',
                    [switch]$Force,
                    [switch]$AllowClobber
                )
                $script:ErrorActionPreference = 'Stop'
                $module = Get-Module -Name $Name -ListAvailable
                if (!$module) {
                    Write-Host "Installing module $Name"
                    Install-Module -Name $Name -Repository $Repository -Force:$Force -AllowClobber:$AllowClobber
                }
            }

            if ($expandC) {
                # Expand C volume as much as possible
                $drive_letter = "C"
                $size = (Get-PartitionSupportedSize -DriveLetter $drive_letter)
                if ($size.SizeMax -gt (Get-Partition -DriveLetter $drive_letter).Size) {
                    echo "Resizing volume"
                    Resize-Partition -DriveLetter $drive_letter -Size $size.SizeMax
                }
            }

            echo "Validate BITS is working"
            $job = Start-BitsTransfer -Source https://aka.ms -Destination $env:TEMP -TransferType Download -Asynchronous
            $count = 0
            while ($job.JobState -ne "Transferred" -and $count -lt 30) {
                if ($job.JobState -eq "TransientError") {
                    throw "BITS transfer failed"
                }
                sleep 6
                $count++
            }
            if ($count -ge 30) {
                throw "BITS transfer failed after 3 minutes. Job state: $job.JobState"
            }

            $creds = [System.Management.Automation.PSCredential]::new($servicePrincipalId, (ConvertTo-SecureString $servicePrincipalSecret -AsPlainText -Force))

            Install-PackageProvider -Name NuGet -MinimumVersion 2.8.5.201 -Force -Confirm:$false

            Install-ModuleIfMissing -Name Az -Repository PSGallery -Force

            Connect-AzAccount -Subscription $subscriptionId -Tenant $tenant -Credential $creds -ServicePrincipal
            echo "login to Azure"

            Install-Module AzSHCI.ARCInstaller -Force -AllowClobber
            Install-Module Az.StackHCI -Force -AllowClobber -RequiredVersion 2.2.3
            Install-Module AzStackHci.EnvironmentChecker -Repository PSGallery -Force -AllowClobber
            Install-ModuleIfMissing Az.Accounts -Force -AllowClobber
            Install-ModuleIfMissing Az.ConnectedMachine -Force -AllowClobber
            Install-ModuleIfMissing Az.Resources -Force -AllowClobber
            echo "Installed modules"
            $id = (Get-AzContext).Tenant.Id
            $token = (Get-AzAccessToken).Token
            $accountid = (Get-AzContext).Account.Id
            Invoke-AzStackHciArcInitialization -SubscriptionID $subscriptionId -ResourceGroup $resourceGroupName -TenantID $id -Region $region -Cloud "AzureCloud" -ArmAccessToken $token -AccountID  $accountid
            $exitCode = $LASTEXITCODE
            $script:ErrorActionPreference = 'Stop'
            if ($exitCode -eq 0) {
                echo "Arc server connected!"
            }
            else {
                throw "Arc server connection failed"
            }

            $ready = $false
            while (!$ready) {
                Connect-AzAccount -Subscription $subscriptionId -Tenant $tenant -Credential $creds -ServicePrincipal
                $extension = Get-AzConnectedMachineExtension -Name "AzureEdgeLifecycleManager" -ResourceGroup $resourceGroupName -MachineName $env:COMPUTERNAME -SubscriptionId $subscriptionId
                if ($extension.ProvisioningState -eq "Succeeded") {
                    $ready = $true
                }
                else {
                    echo "Waiting for LCM extension to be ready"
                    Start-Sleep -Seconds 30
                }
            }

        } -ArgumentList $subscriptionId, $resourceGroupName, $region, $tenant, $servicePrincipalId, $servicePrincipalSecret
        break
    }
    catch {
        echo "Error in retry ${count}:`n$_"
        sleep 600
    }
    finally {
        if ($session) {
            Remove-PSSession -Session $session
        }
    }
}

if ($count -ge 6) {
    throw "Failed to connect Arc server after 6 retries."
}
