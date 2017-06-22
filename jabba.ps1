<#
.SYNOPSIS
A powershell wrapper for the Jabba JVM manager, commonly aliased as "jabba".
.DESCRIPTION
Jabba makes it possible to manage multiple JVMs on your system; you can list available JVMs, install remotely-available JVMs and most importantly, select a JVM to work with. Selecting a JVM in this case means manipulating the PATH and JAVA_HOME variables.

.EXAMPLE
Invoke-Jabba help

Show the help output for the jabba executable.

.EXAMPLE
Invoke-Jabba ls

List currently known and available JVMs.

.EXAMPLE
Invoke-Jabba link system@1.8.131 """C:\Program Files\Java\oracle-jdk8"""

Tells Jabba about an already-installed JVM.

.EXAMPLE
Invoke-Jabba use system@1.8.131

Tells jabba to activate the system installed JDK version 1.8.131.
#>
function Invoke-Jabba
{
    $env:JABBA_HOME=$env:USERPROFILE + "\.jabba"
    $fd3=$([System.IO.Path]::GetTempFileName())
    $command=$env:JABBA_HOME + "\bin\jabba.exe $args --fd3 $fd3"
    & { $env:JABBA_SHELL_INTEGRATION="ON"; Invoke-Expression $command }
    $fd3content=$(Get-Content $fd3)
    if ($fd3content) {
        $expression=$fd3content.replace("export ","`$env:") -join "`n"
        if (-not $expression -eq "") { Invoke-Expression $expression }
    }
    Remove-Item -Force $fd3
}
