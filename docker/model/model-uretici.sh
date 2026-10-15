#!/usr/bin/env bash
# Cevap üreten küçük dil modelini (Qwen3.5-0.8B) llama.cpp sunucusu olarak başlatır.
# Bonus (ANSWER komutu) içindir. Koordinatör bunu çocuk process olarak başlatır
# ve 127.0.0.1 üzerinden HTTP ile konuşur.
#
# Kullanım: model-uretici [port]          (varsayılan port 8081)
#
# Hazır olduğunu anlamak için:  curl -s 127.0.0.1:8081/health   -> {"status":"ok"}
# Soru sormak için (OpenAI uyumlu):
#   curl -s 127.0.0.1:8081/v1/chat/completions -d '{"messages":[{"role":"user","content":"Merhaba"}]}'
#
# Ortam değişkenleri:
#   URETICI_THREADS   iş parçacığı sayısı (varsayılan 4)
#   URETICI_MODEL     GGUF dosyası (varsayılan /opt/modeller/qwen3.5-0.8b-q8_0.gguf)
set -euo pipefail

port=${1:-8081}
export LD_LIBRARY_PATH=/opt/llama${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}
exec /opt/llama/llama-server \
    --model "${URETICI_MODEL:-/opt/modeller/qwen3.5-0.8b-q8_0.gguf}" \
    --host 127.0.0.1 --port "$port" \
    --threads "${URETICI_THREADS:-4}" \
    --ctx-size 4096 \
    --reasoning off \
    --no-webui
