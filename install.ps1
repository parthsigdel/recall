$ErrorActionPreference = "Stop"

$Repo = "parthsigdel/recall"
$BinName = "recall"
$InstallDir = "$env:LOCALAPPDATA\Programs\recall"

# detect arch
$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
    "AMD64" { $Arch = "amd64" }
    "ARM64" { $Arch = "arm64" }
    default { Write-Error "Unsupported architecture: $Arch"; exit 1 }
}

$Archive = "${BinName}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/latest/download/$Archive"

$TmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP ([System.Guid]::NewGuid()))
$ZipPath = Join-Path $TmpDir $Archive

Write-Host "Downloading $Url..."
Invoke-WebRequest -Uri $Url -OutFile $ZipPath

Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Move-Item -Force (Join-Path $TmpDir "$BinName.exe") (Join-Path $InstallDir "$BinName.exe")

Remove-Item -Recurse -Force $TmpDir

# add to user PATH if not already present
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to your PATH. Restart your terminal for it to take effect."
}

Write-Host "Installed $BinName to $InstallDir\$BinName.exe"
