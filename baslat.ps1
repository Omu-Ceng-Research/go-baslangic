# Sistem Programlama Go ortamını Docker'da başlatır (Windows).
#
#   baslat.cmd dosyasına çift tıklayın ya da PowerShell'de:
#   powershell -ExecutionPolicy Bypass -File .\baslat.ps1 [temel|model|llm]
#
#   model: gerçek embedding modeli olan imaj (koordinatör için)
#   llm:   model + cevap üreten dil modeli (bonus için)
#
# Konteyner zaten açıksa ona ikinci bir terminal açar.

param([ValidateSet("temel", "model", "llm")][string]$Tur = "temel")

$Etiket = if ($Tur -eq "temel") { "latest" } else { $Tur }
$Image = if ($env:SYSPROG_IMAGE) { $env:SYSPROG_IMAGE } else { "ghcr.io/omu-ceng-research/sysprog-go:$Etiket" }
$Ad = "sysprog"
Set-Location $PSScriptRoot

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "Docker kurulu değil. README.md'deki kurulum adımlarına bakın." -ForegroundColor Red
    exit 1
}

docker info *> $null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Docker çalışmıyor. Docker Desktop'ı başlatın ve tekrar deneyin." -ForegroundColor Red
    exit 1
}

$acik = docker ps -q -f "name=^$Ad$"
if ($acik) {
    docker exec -it -u 0 -w /calisma $Ad bash
    exit $LASTEXITCODE
}

docker image inspect $Image *> $null
if ($LASTEXITCODE -ne 0) {
    Write-Host "İmaj indiriliyor: $Image"
    docker pull $Image
    if ($LASTEXITCODE -ne 0) {
        Write-Host "İndirilemedi, yerelde derleniyor (birkaç dakika sürebilir)..."
        docker build -f docker/Dockerfile --target $Tur -t $Image .
        if ($LASTEXITCODE -ne 0) { exit 1 }
    }
}

docker run -it --rm --name $Ad --hostname sysprog `
    --cap-add SYS_PTRACE `
    -v "${PWD}:/kod" `
    -v sysprog-calisma:/calisma `
    $Image
