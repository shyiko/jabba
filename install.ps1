$ErrorActionPreference = "Stop"

$jabbaHome = if ($env:JABBA_HOME) { $env:JABBA_HOME } else { if ($env:JABBA_DIR) { $env:JABBA_DIR } else { "$env:USERPROFILE\.jabba" } }
$jabbaVersion = if ($env:JABBA_VERSION) { $env:JABBA_VERSION } else { "latest" }

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
    Invoke-WebRequest https://github.com/shyiko/jabba/releases/download/$jabbaVersion/jabba-$jabbaVersion-windows-amd64.exe -UseBasicParsing -OutFile $jabbaHome/bin/jabba.exe
}

$ErrorActionPreference="SilentlyContinue"
& $jabbaHome\bin\jabba.exe --version | Out-Null
$binaryValid = $?
$ErrorActionPreference="Continue"
if (-not $binaryValid)
{
    Write-Host @"
$jabbaHome\bin\jabba does not appear to be a valid binary.

Check your Internet connection / proxy settings and try again.
if the problem persists - please create a ticket at https://github.com/shyiko/jabba/issues.
"@
    exit 1
}

@"
`$env:JABBA_HOME="$jabbaHome"
if (Test-Path "`$env:JABBA_HOME/jdk/default") {
    `$env:JAVA_HOME = "`$env:JABBA_HOME\jdk\default"
    `$env:Path = "`$env:JAVA_HOME\bin;`$env:Path"
}

function jabba
{
    `$fd3=`$([System.IO.Path]::GetTempFileName())
    `$command="& '$jabbaHome\bin\jabba.exe' `$args --fd3 ```"`$fd3```""
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

Installation completed
(if you have any problems please report them at https://github.com/shyiko/jabba/issues)
"@
