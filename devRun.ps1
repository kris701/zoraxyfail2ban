if (!(Test-Path "./dev")){
    New-Item -ItemType Directory -Path "./dev"
    cd dev
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest https://github.com/tobychui/zoraxy/releases/latest/download/zoraxy_linux_amd64 -OutFile zoraxy
    cd ..
}

if (!(Test-Path "./dev/plugins/zoraxyfail2ban")) {
    New-Item -ItemType Directory -Path "./dev/plugins/zoraxyfail2ban"
}

$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
go build

$currentDir = (Get-Item .).FullName
wsl --user root killall -s 9 zoraxy
Copy-Item "zoraxyfail2ban" -Destination "./dev/plugins/zoraxyfail2ban/"
$j = Start-Job -Name ZoraxyServer -ArgumentList $currentDir -ScriptBlock {
    Set-Location $args[0]
    $currentDir = (Get-Item .).FullName
    echo $currentDir
    wsl --cd "$currentDir/dev" --user root ./zoraxy -noauth=true -port=:8564 -sshlb=true
}
echo "Done!"
