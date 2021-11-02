$ErrorActionPreference = "Stop"

$sep = [IO.Path]::DirectorySeparatorChar
$jabbaHome = if ($env:JABBA_HOME) { 
    $env:JABBA_HOME 
} else { 
    if ($env:JABBA_DIR) { 
        $env:JABBA_DIR 
    } else {
        if($env:USERPROFILE){
            "$env:USERPROFILE" + $sep + ".jabba" 
        }else{
            "$env:HOME" + $sep + ".jabba"
        }
    } 
}
$jabbaVersion = if ($env:JABBA_VERSION) { $env:JABBA_VERSION } else { "latest" }
# The Windows values of the Platform enum are:
# 0 (Win32NT), 1 (Win32S), 2 (Win32Windows) and 3 (WinCE).
# Other values are larger and indicate non Windows operating systems
$isOnWindows = [System.Environment]::OSVersion.Platform.value__ > 3
$jabbaExecutableName = $isOnWindows ? "jabba.exe" : "jabba"

if ($jabbaVersion -eq "latest")
{
    # resolving "latest" to an actual tag
    $jabbaVersion = (Invoke-RestMethod https://api.github.com/repos/shyiko/jabba/releases/latest).body
}

if ($jabbaVersion -notmatch '^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.+-]+)?$')
{
    Write-Host "'$jabbaVersion' is not a valid version."
    exit 1
}

Write-Host "Installing v$jabbaVersion...`n"

New-Item -Type Directory -Force $jabbaHome/bin | Out-Null

if ($env:JABBA_MAKE_INSTALL -eq "true")
{
    Copy-Item jabba.exe $jabbaHome/bin
}
else
{
    # $isOnWindows, see top of the file
    # MacOSX enum value: 4
    if($isOnWindows){
        Invoke-WebRequest https://github.com/shyiko/jabba/releases/download/$jabbaVersion/jabba-$jabbaVersion-windows-amd64.exe -UseBasicParsing -OutFile $jabbaHome/bin/$jabbaExecutableName
    }
    elseif([System.Environment]::OSVersion.Platform.value__ -eq 4){
        Invoke-WebRequest https://github.com/shyiko/jabba/releases/download/$jabbaVersion/jabba-$jabbaVersion-darwin-amd64 -UseBasicParsing -OutFile $jabbaHome/bin/$jabbaExecutableName
    }else{
        $osArch = [System.Environment]::Is64BitOperatingSystem ? "amd64" : "386"
        Invoke-WebRequest https://github.com/shyiko/jabba/releases/download/$jabbaVersion/jabba-$jabbaVersion-linux-${OSARCH} -UseBasicParsing -OutFile $jabbaHome/bin/$jabbaExecutableName
    }
}

$ErrorActionPreference="SilentlyContinue"

if($isOnWindows){
    & "$jabbaHome\bin\$jabbaExecutableName" --version | Out-Null
}else{
	chmod a+x "$jabbaHome/bin/$jabbaExecutableName"
    & "$jabbaHome/bin/$jabbaExecutableName" --version | Out-Null
}

$binaryValid = $?
$ErrorActionPreference="Continue"
if (-not $binaryValid)
{
    Write-Host -ForegroundColor Yellow @"
$jabbaHome\bin\$jabbaExecutableName does not appear to be a valid binary.

Check your Internet connection / proxy settings and try again.
if the problem persists - please create a ticket at https://github.com/shyiko/jabba/issues.
"@
    exit 1
}

@"
`$env:JABBA_HOME="$jabbaHome"

function jabba
{
    `$fd3=`$([System.IO.Path]::GetTempFileName())
    `$command="& '$jabbaHome\bin\$jabbaExecutableName' `$args --fd3 ```"`$fd3```""
    & { `$env:JABBA_SHELL_INTEGRATION="ON"; Invoke-Expression `$command }
    `$fd3content=`$(Get-Content `$fd3)
    if (`$fd3content) {
        `$expression=`$fd3content.replace("export ","```$env:").replace("unset ","Remove-Item env:") -join "``n"
        if (-not `$expression -eq "") { Invoke-Expression `$expression }
    }
    Remove-Item -Force `$fd3
}
"@ | Out-File $jabbaHome/jabba.ps1

$sourceJabba="if (Test-Path `"$jabbaHome\jabba.ps1`") { . `"$jabbaHome\jabba.ps1`" }"

if (-not $(Test-Path $profile))
{
    New-Item -Path $profile -Type File -Force | Out-Null
}

if ("$(Get-Content $profile | Select-String "\\jabba.ps1")" -eq "")
{
    Write-Host "Adding source string to $profile"
    "`n$sourceJabba`n" | Out-File -Append -Encoding ASCII $profile
}
else
{
    Write-Host "Skipped update of $profile (source string already present)"
}

. $jabbaHome\jabba.ps1

Write-Host @"

Installation completed (you might need to restart your terminal for the jabba command to be available)
(if you have any problems please report them at https://github.com/shyiko/jabba/issues)
"@
