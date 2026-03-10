[CmdletBinding()]
param(
    [ValidateSet('auto', 'local', 'release')]
    [string]$Source = 'auto',

    [string]$Version = 'latest',

    [string]$InstallDir = (Join-Path $env:LOCALAPPDATA 'Programs\Codex Themes\bin'),

    [string]$Repo = 'ychampion/codex-themes',

    [string]$Theme,

    [string]$StatusLine,

    [switch]$NoPathPersist,

    [switch]$NoPrompt,

    [switch]$RunSetup
)

$ErrorActionPreference = 'Stop'

function Write-Step {
    param([string]$Message)
    Write-Host "[codex-theme] $Message"
}

function Test-Command {
    param([string]$Name)
    return $null -ne (Get-Command $Name -ErrorAction SilentlyContinue)
}

function Get-RepoRoot {
    if ([string]::IsNullOrWhiteSpace($PSScriptRoot)) {
        return $null
    }

    $candidate = Split-Path -Parent $PSScriptRoot
    if ((Test-Path (Join-Path $candidate 'go.mod')) -and (Test-Path (Join-Path $candidate 'cmd\codex-theme\main.go'))) {
        return $candidate
    }

    return $null
}

function Ensure-InstallDir {
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
}

function Add-ToPath {
    $currentSessionParts = @($env:PATH -split ';' | Where-Object { $_ -and $_.Trim() -ne '' })
    if ($currentSessionParts -notcontains $InstallDir) {
        $env:PATH = "$InstallDir;$env:PATH"
    }

    if ($NoPathPersist) {
        return
    }

    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    $userParts = @($userPath -split ';' | Where-Object { $_ -and $_.Trim() -ne '' })
    if ($userParts -notcontains $InstallDir) {
        $newUserPath = if ([string]::IsNullOrWhiteSpace($userPath)) {
            $InstallDir
        } else {
            "$userPath;$InstallDir"
        }
        [Environment]::SetEnvironmentVariable('Path', $newUserPath, 'User')
        Write-Step "Added $InstallDir to your user PATH."
    }
}

function Build-LocalBinary {
    param([string]$RepoRoot)

    Ensure-InstallDir
    $targetPath = Join-Path $InstallDir 'codex-theme.exe'
    $repoBinary = Join-Path $RepoRoot 'codex-theme.exe'

    if (Test-Command 'go') {
        Write-Step 'Building local binary with Go.'
        & go build -o $targetPath ./cmd/codex-theme
        return $targetPath
    }

    if (Test-Path $repoBinary) {
        Write-Step 'Copying prebuilt local binary.'
        Copy-Item $repoBinary $targetPath -Force
        return $targetPath
    }

    throw 'Local install requested, but Go is unavailable and no prebuilt codex-theme.exe was found in the repo root.'
}

function Get-PlatformAssetName {
    $arch = switch ($env:PROCESSOR_ARCHITECTURE) {
        'ARM64' { 'arm64' }
        default { 'x64' }
    }

    return "codex-theme_windows_${arch}.zip"
}

function Download-File {
    param(
        [string]$Url,
        [string]$Destination
    )

    Write-Step "Downloading $Url"
    Invoke-WebRequest -Uri $Url -OutFile $Destination -UseBasicParsing
}

function Install-ReleaseBinary {
    Ensure-InstallDir

    $asset = Get-PlatformAssetName
    $tag = $Version

    if ($Version -eq 'latest') {
        Write-Step 'Resolving latest GitHub release.'
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Headers @{ 'User-Agent' = 'codex-theme-installer' }
        $tag = $release.tag_name
    } elseif (-not $Version.StartsWith('v')) {
        $tag = "v$Version"
    }

    $zipPath = Join-Path ([System.IO.Path]::GetTempPath()) ("codex-theme-$tag-$([guid]::NewGuid()).zip")
    $extractDir = Join-Path ([System.IO.Path]::GetTempPath()) ("codex-theme-$tag-$([guid]::NewGuid())")
    $url = "https://github.com/$Repo/releases/download/$tag/$asset"

    Download-File -Url $url -Destination $zipPath
    Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

    $binary = Get-ChildItem -Path $extractDir -Filter 'codex-theme.exe' -Recurse | Select-Object -First 1
    if (-not $binary) {
        throw "Downloaded archive $asset did not contain codex-theme.exe"
    }

    $targetPath = Join-Path $InstallDir 'codex-theme.exe'
    Copy-Item $binary.FullName $targetPath -Force
    Remove-Item $zipPath -Force -ErrorAction SilentlyContinue
    Remove-Item $extractDir -Recurse -Force -ErrorAction SilentlyContinue
    return $targetPath
}

function Resolve-InstallMode {
    $repoRoot = Get-RepoRoot

    switch ($Source) {
        'local' { return @{ Mode = 'local'; RepoRoot = $repoRoot } }
        'release' { return @{ Mode = 'release'; RepoRoot = $repoRoot } }
        default {
            if ($repoRoot) {
                return @{ Mode = 'local'; RepoRoot = $repoRoot }
            }
            return @{ Mode = 'release'; RepoRoot = $repoRoot }
        }
    }
}

function Should-RunSetupInteractively {
    if (-not [string]::IsNullOrWhiteSpace($Theme)) {
        return $false
    }

    if ($RunSetup) {
        return $true
    }

    if ($NoPrompt) {
        return $false
    }

    if (-not $Host.UI) {
        return $false
    }

    $answer = Read-Host 'Launch the guided theme setup now? [Y/n]'
    return [string]::IsNullOrWhiteSpace($answer) -or $answer.Trim().ToLowerInvariant() -in @('y', 'yes')
}

$modeInfo = Resolve-InstallMode
Write-Step "Install mode: $($modeInfo.Mode)"

$installedBinary = if ($modeInfo.Mode -eq 'local') {
    if (-not $modeInfo.RepoRoot) {
        throw 'Could not detect the repo root for a local install.'
    }

    Push-Location $modeInfo.RepoRoot
    try {
        Build-LocalBinary -RepoRoot $modeInfo.RepoRoot
    }
    finally {
        Pop-Location
    }
} else {
    Install-ReleaseBinary
}

Add-ToPath

Write-Step "Installed to $installedBinary"
Write-Step 'You can now run codex-theme from this PowerShell window.'

if (-not [string]::IsNullOrWhiteSpace($Theme)) {
    $arguments = @('apply', $Theme)
    if (-not [string]::IsNullOrWhiteSpace($StatusLine)) {
        $arguments += @('--status-line', $StatusLine)
    }
    $message = "Applying theme $Theme"
    if (-not [string]::IsNullOrWhiteSpace($StatusLine)) {
        $message += " with status line $StatusLine"
    }
    Write-Step $message
    & $installedBinary @arguments
    if ($LASTEXITCODE -eq 0) {
        Write-Step 'Done. Restart Codex to see the new footer/status line.'
    }
    exit $LASTEXITCODE
}

if (Should-RunSetupInteractively) {
    Write-Step 'Opening guided setup.'
    & $installedBinary setup
    exit $LASTEXITCODE
}

Write-Step 'Next command: codex-theme setup'
