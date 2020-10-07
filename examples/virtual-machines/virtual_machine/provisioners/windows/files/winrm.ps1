$profiles = Get-NetConnectionProfile
Foreach ($i in $profiles) {
    Write-Host ("Updating Interface ID {0} to be Private.." -f $profiles.InterfaceIndex)
    Set-NetConnectionProfile -InterfaceIndex $profiles.InterfaceIndex -NetworkCategory Private
}

Write-Host "Obtaining the Thumbprint of the Certificate from KeyVault"
$Thumbprint = (Get-ChildItem -Path Cert:\LocalMachine\My | Where-Object {$_.Subject -match "$ComputerName"}).Thumbprint

Write-Host "Enable HTTPS in WinRM.."
winrm create winrm/config/Listener?Address=*+Transport=HTTPS "@{Hostname=`"$ComputerName`"; CertificateThumbprint=`"$Thumbprint`"}"

Write-Host "Enabling Basic Authentication.."
winrm set winrm/config/service/Auth "@{Basic=`"true`"}"

Write-Host "Re-starting the WinRM Service"
net stop winrm
net start winrm

Write-Host "Open Firewall Ports"
netsh advfirewall firewall add rule name="Windows Remote Management (HTTPS-In)" dir=in action=allow protocol=TCP localport=5986
