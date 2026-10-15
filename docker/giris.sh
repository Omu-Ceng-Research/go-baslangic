#!/usr/bin/env bash
# Konteyner giriş scripti: /kod boşsa dersleri kopyalar, karşılama mesajı basar.
set -e

if [[ ! -f /kod/go.mod ]]; then
    if [[ -z "$(ls -A /kod 2>/dev/null)" ]]; then
        cp -r /opt/go-baslangic/. /kod/
        echo "Dersler /kod klasörüne kopyalandı (bilgisayarınızdaki klasörde de görünür)."
    else
        echo "UYARI: /kod boş değil ama go.mod yok. Dersler kopyalanmadı." >&2
        echo "Boş bir klasörden ya da go-baslangic klasörünün içinden başlatın." >&2
    fi
fi

if [[ -t 1 && "$1" == "bash" ]]; then
    cat <<EOF

  Sistem Programlama — Go ortamı ($(go version | cut -d' ' -f3), $(uname -m))

    ders 01                 01. dersi derle ve çalıştır
    ders 10 -adet 300000    argüman verilebilir
    cd /calisma             derlenen programlar ve üretilen dosyalar burada
    strace -c ./d10         sistem çağrılarını say
    man 2 mmap              sistem çağrısının belgesi

  İkinci terminal: başlatma scriptini (baslat.sh / baslat.cmd) tekrar çalıştırın.

EOF
fi

exec "$@"
