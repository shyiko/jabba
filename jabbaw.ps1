$ErrorActionPreference = "Stop"

$jabbaHome = if ($env:JABBA_HOME) { $env:JABBA_HOME } else { if ($env:JABBA_DIR) { $env:JABBA_DIR } else { "$env:USERPROFILE\.jabba" } }
$jabbaVersion = if ($env:JABBA_VERSION) { $env:JABBA_VERSION } else { "latest" }

if ($jabbaVersion -eq "latest") {
  # resolving "latest" to an actual tag
  $jabbaVersion = (Invoke-RestMethod https://api.github.com/repos/Jabba-Team/jabba/releases/latest).body
}

$ErrorActionPreference = "SilentlyContinue"
& $jabbaHome\bin\jabba.exe --version | Out-Null
$binaryValid = $?
$ErrorActionPreference = "Continue"

if ($binaryValid) {
  $realVersion = Invoke-Expression "$jabbaHome\bin\jabba.exe --version"
}

if (-not $binaryValid -or $jabbaVersion -ne $realVersion) {
    
  [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
  Invoke-Expression (
    Invoke-WebRequest https://github.com/Jabba-Team/jabba/raw/main/install.ps1 -UseBasicParsing
  ).Content

}

Invoke-Expression "$jabbaHome\bin\jabba.exe install"
$env:JAVA_HOME = Invoke-Expression "$jabbaHome\bin\jabba.exe which --home"
Invoke-Expression "& $args" 
