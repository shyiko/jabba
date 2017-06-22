@echo off

:: Set environment variables.
set "JABBA_HOME=%USERPROFILE%\.jabba"
set "TEMP_ENV=%TEMP%\jabba.env"
set "TEMP_REP=%TEMP%\jabba.rep"
set "TEMP_BAT=%TEMP%\jabba.bat"

:: Run jabba, write to TEMP_ENV.
"%JABBA_HOME%\bin\jabba.exe" %1 %2 %3 %4 %5 --fd3 "%TEMP_ENV%"

if exist "%TEMP_ENV%" (
    :: Convert from powershell to batch environment schript.
    powershell -Command "(gc %TEMP_ENV%) -replace 'export ', 'set \"' ^| Out-File -encoding ASCII %TEMP_REP%"
    powershell -Command "(gc %TEMP_REP%) -replace '=\"', '=' ^| Out-File -encoding ASCII %TEMP_BAT%"

    :: Exectute within current shell (i.e, no call).
    "%TEMP_BAT%"

    :: Clean up files.
    del "%TEMP_ENV%" "%TEMP_REP%" "%TEMP_BAT%"
)

:: Unset used environment variables.
set JABBA_HOME=
set TEMP_ENV=
set TEMP_REP=
set TEMP_BAT=
