# Enable TLS 1.2 for secure downloads
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# Function to check if a port is available
function Test-PortAvailable {
    param ([int]$port)
    $tcpListener = New-Object System.Net.Sockets.TcpListener([System.Net.IPAddress]::Any, $port)
    try {
        $tcpListener.Start()
        $tcpListener.Stop()
        return $true
    } catch {
        return $false
    }
}

# Install Chocolatey if not already installed
if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Chocolatey..."
    Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
}

# Install JDK 8 using Chocolatey
if (-not (Get-Command java -ErrorAction SilentlyContinue)) {
    Write-Host "Installing JDK 8..."
    choco install -y jdk8
} else {
    Write-Host "JDK already installed."
}

# Set JAVA_HOME environment variable (required for Tomcat)
$javaPath = (Get-Command java).Path
$javaHome = [System.IO.Path]::GetDirectoryName($javaPath)
$env:JAVA_HOME = $javaHome
[System.Environment]::SetEnvironmentVariable("JAVA_HOME", $javaHome, "Machine")

# Install XAMPP using Chocolatey
if (-not (Test-Path "C:\xampp")) {
    Write-Host "Installing XAMPP..."
    choco install -y xampp
} else {
    Write-Host "XAMPP already installed."
}

# Check if Tomcat port is available (default is 8080)
$port = 8080
if (-not (Test-PortAvailable $port)) {
    Write-Host "Port 8080 is not available. Changing Tomcat to use port 8081."
    $port = 8081

    # Update the Tomcat server.xml to use the new port
    $serverXmlPath = "C:\xampp\tomcat\conf\server.xml"
    (Get-Content $serverXmlPath) -replace 'port="8080"', 'port="8081"' | Set-Content $serverXmlPath
}

# Start XAMPP Control Panel
Write-Host "Starting XAMPP services..."
Start-Process "C:\xampp\xampp-control.exe"

# Start Tomcat via XAMPP
Write-Host "Starting Tomcat..."
Start-Process "C:\xampp\tomcat_start.bat"

Write-Host "XAMPP and Tomcat setup completed."


