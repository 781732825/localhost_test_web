# Generate self-signed certificate and key using OpenSSL

# Create certificate directory
$certDir = "certs"
if (-not (Test-Path $certDir)) {
    New-Item -ItemType Directory -Path $certDir | Out-Null
    Write-Host "Creating certificate directory: $certDir"
}

# Certificate and key file paths
$certFile = "$certDir\server.crt"
$keyFile = "$certDir\server.key"

# Check if OpenSSL is available
$openssl = $null
try {
    $openssl = Get-Command openssl -ErrorAction Stop
    Write-Host "OpenSSL found: $($openssl.Source)"
} catch {
    Write-Host "OpenSSL not found in PATH. Please install OpenSSL and add it to your PATH." -ForegroundColor Red
    exit 1
}

# Generate private key
Write-Host "Generating private key..."
& openssl genrsa -out $keyFile 2048
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to generate private key." -ForegroundColor Red
    exit 1
}
Write-Host "Private key generated: $keyFile" -ForegroundColor Green

# Generate self-signed certificate
Write-Host "Generating self-signed certificate..."
& openssl req -new -x509 -key $keyFile -out $certFile -days 365 -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost,DNS:example.com,IP:127.0.0.1"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to generate certificate." -ForegroundColor Red
    exit 1
}
Write-Host "Certificate generated: $certFile" -ForegroundColor Green

# Update configuration file
$configFile = "config.json"
$config = Get-Content $configFile | ConvertFrom-Json

# Update HTTPS configuration for port 443
foreach ($port in $config.ports) {
    if ($port.port -eq 443) {
        $port.https = $true
        $port | Add-Member -NotePropertyName "cert" -NotePropertyValue $certFile -Force
        $port | Add-Member -NotePropertyName "key" -NotePropertyValue $keyFile -Force
        
        # Remove PFX properties if they exist
        if ($port.PSObject.Properties['pfx']) {
            $port.PSObject.Properties.Remove('pfx')
        }
        if ($port.PSObject.Properties['pfx_password']) {
            $port.PSObject.Properties.Remove('pfx_password')
        }
    }
}

# Save updated configuration
$config | ConvertTo-Json -Depth 10 | Set-Content $configFile
Write-Host "Configuration file updated, HTTPS enabled" -ForegroundColor Green

Write-Host "\nNOTE: Since we're using a self-signed certificate, browsers will show a security warning." -ForegroundColor Yellow
Write-Host "You'll need to accept the risk to proceed to the website." -ForegroundColor Yellow