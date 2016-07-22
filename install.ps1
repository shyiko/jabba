$ErrorActionPreference = "Stop"

$jabbaDir = If ($env:JABBA_DIR) { $env:JABBA_DIR } else { "$env:USERPROFILE\.jabba" }
$jabbaVersion = If ($env:JABBA_VERSION) { $env:JABBA_VERSION } else { "latest" }

If ($jabbaVersion -eq "latest")
{
    # resolving "latest" to an actual tag
    $jabbaVersion = [System.Text.Encoding]::UTF8.GetString((wget https://shyiko.github.com/jabba/latest -UseBasicParsing).Content).Trim()
}

If ($jabbaVersion -notmatch '^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.+-]+)?$')
{
    echo "'$jabbaVersion' is not a valid version."
    exit 1
}

echo "Installing v$jabbaVersion..."
echo ""

mkdir -Force $jabbaDir/bin | Out-Null

(wget https://github.com/shyiko/jabba/releases/download/$jabbaVersion/jabba-$jabbaVersion-windows-amd64.exe -UseBasicParsing).Content > $jabbaDir/bin/jabba.exe

$ErrorActionPreference="SilentlyContinue"
& $jabbaDir\bin\jabba.exe --version | Out-Null
$binaryValid = $?
$ErrorActionPreference="Continue"
if (-not $binaryValid)
{
    echo "$jabbaDir\bin\jabba does not appear to be a valid binary.

Check your Internet connection / proxy settings and try again.
If the problem persists - please create a ticket at https://github.com/shyiko/jabba/issue."
    exit 1
}

echo @"
function jabba
{
    `$fd3=`$([System.IO.Path]::GetTempFileName())
    `$command="$jabbaDir\bin\jabba.exe `$args --fd3 `$fd3"
    & { `$env:JABBA_SHELL_INTEGRATION="ON"; Invoke-Expression `$command }
    `$fd3content=`$(cat `$fd3)
    if (`$fd3content) {
        `$expression=`$fd3content.replace("export ","```$env:") -join "``n"
        if (-not `$expression -eq "") { Invoke-Expression `$expression }
    }
    rm -Force `$fd3
}
"@ > $jabbaDir/jabba.ps1

$sourceJabba="if (Test-Path `"$jabbaDir\jabba.ps1`") { . `"$jabbaDir\jabba.ps1`" }"

if (-not $(Test-Path $profile))
{
    New-Item -path $profile -type file –force | Out-Null
}

if ("$(cat $profile | Select-String "\\jabba.ps1")" -eq "")
{
    echo "Adding source string to $profile"
    echo "`n$sourceJabba`n" >> "$profile"
}
else
{
    echo "Skipped update of $profile (source string already present)"
}

echo ""
echo "Installation completed`
(if you have any problems please report them at https://github.com/shyiko/jabba/issue)"
