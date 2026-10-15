#!/usr/bin/env bash
# Kullanım: ders <numara|klasör> [programın argümanları...]
#   ders 10
#   ders 10-mmap -adet 300000
#   ders 15            (test + benchmark)
#
# Dersi /kod'dan derler, binary'yi /calisma/dNN olarak koyar ve /calisma'da
# çalıştırır. Böylece üretilen dosyalar hızlı Linux dosya sisteminde kalır ve
# aynı binary'yi strace ile tekrar çalıştırabilirsiniz.
set -euo pipefail

KOD=/kod
CALISMA=${CALISMA:-/calisma}

if [[ $# -lt 1 ]]; then
    echo "Kullanım: ders <numara> [argümanlar...]   örn: ders 10 -adet 300000" >&2
    echo "Dersler:" >&2
    (cd "$KOD" && ls -d [0-9][0-9]-*/ | tr -d / | sed 's/^/  /') >&2
    exit 1
fi

istek=$1
shift
no=${istek:0:2}
klasor=$(cd "$KOD" && ls -d "${no}"-*/ 2>/dev/null | head -1 | tr -d /)
if [[ -z "$klasor" ]]; then
    echo "Ders bulunamadı: $istek" >&2
    exit 1
fi

mkdir -p "$CALISMA"
cd "$KOD"

# 15 bir kütüphane paketidir: çalıştırılmaz, test edilir.
if ! grep -qs '^package main' "$klasor"/*.go; then
    echo ">> go test -v ./$klasor && go test -bench . -benchmem ./$klasor"
    go test -v "./$klasor"
    go test -run '^$' -bench . -benchmem "./$klasor" "$@"
    exit
fi

bin="$CALISMA/d$no"
echo ">> go build -o $bin ./$klasor"
go build -o "$bin" "./$klasor"
echo ">> cd $CALISMA && ./d$no $*"
echo
cd "$CALISMA"
exec "./d$no" "$@"
