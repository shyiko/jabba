$ErrorActionPreference = 'Stop';

$packageName  = 'jabba'
$installType  = 'exe'
$url          = 'https://github.com/shyiko/jabba/releases/download/$version$/jabba-$version$-windows-amd64.exe'
$checksum     = '$checksum$'
$checksumType = 'sha256'
$silentArgs   = '/VERYSILENT'


# Install
Install-ChocolateyPackage $name $installType $silentArgs $url

