#!/usr/bin/env bash
# Sistem Programlama Go ortamını Docker'da başlatır (macOS ve Linux).
#
#   ./baslat.sh          ilk çalıştırma: konteyneri açar
#   ./baslat.sh          konteyner zaten açıksa: ona ikinci bir terminal açar
#   ./baslat.sh model    gerçek embedding modeli olan imajla (koordinatör için)
#   ./baslat.sh llm      model + cevap üreten dil modeli (bonus için)
#
# Bu klasör konteynerde /kod olarak görünür; dosyaları kendi editörünüzle
# düzenleyebilirsiniz. Derlenen programlar /calisma'ya yazılır.
set -euo pipefail

TUR=${1:-temel}
case "$TUR" in
    temel) ETIKET=latest ;;
    model|llm) ETIKET=$TUR ;;
    *) echo "Kullanım: $0 [temel|model|llm]" >&2; exit 1 ;;
esac
IMAGE=${SYSPROG_IMAGE:-ghcr.io/omu-ceng-research/sysprog-go:$ETIKET}
AD=sysprog
cd "$(dirname "$0")"

if ! command -v docker >/dev/null; then
    echo "Docker kurulu değil. README.md'deki kurulum adımlarına bakın." >&2
    exit 1
fi
if ! docker info >/dev/null 2>&1; then
    echo "Docker çalışmıyor. Docker Desktop'ı (ya da OrbStack/Colima'yı) başlatın." >&2
    exit 1
fi

# Konteyner açıksa ikinci terminal. Root olarak açılır ki strace -p ile
# ilk terminaldeki programa bağlanılabilsin.
if [[ -n "$(docker ps -q -f "name=^${AD}$")" ]]; then
    exec docker exec -it -u 0 -w /calisma "$AD" bash
fi

if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
    echo "İmaj indiriliyor: $IMAGE"
    if ! docker pull "$IMAGE"; then
        echo "İndirilemedi, yerelde derleniyor (birkaç dakika sürebilir)..."
        docker build -f docker/Dockerfile --target "$TUR" -t "$IMAGE" .
    fi
fi

args=(run -it --rm --name "$AD" --hostname sysprog
      --cap-add SYS_PTRACE
      -v "$PWD":/kod
      -v sysprog-calisma:/calisma)

# Linux'ta oluşturulan dosyalar root'a değil size ait olsun.
if [[ "$(uname -s)" == Linux ]]; then
    args+=(--user "$(id -u):$(id -g)")
fi

exec docker "${args[@]}" "$IMAGE"
