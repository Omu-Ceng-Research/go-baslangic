# Sistem Programlama — Go çalışma ortamı
#
# Windows, macOS ve Linux'ta aynı Linux ortamını sağlar. Üç katman (target) var:
#
#   temel  Go, strace, gcc, netcat, man sayfaları. Dersler ve depolama düğümleri.
#   model  temel + gerçek embedding modeli (multilingual-e5-small). Koordinatör.
#   llm    model + cevap üreten dil modeli (Qwen3.5-0.8B, llama.cpp). Bonus.
#
# Derleme (go-baslangic klasöründen):
#   docker build -f docker/Dockerfile --target temel -t sysprog-go:latest .
#   docker build -f docker/Dockerfile --target model -t sysprog-go:model .
#   docker build -f docker/Dockerfile --target llm   -t sysprog-go:llm .

# ---------------------------------------------------------------------------
FROM golang:1.27-trixie AS temel

# man sayfaları sistem çağrılarını öğrenmek için önemli (man 2 open, man 2 mmap).
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        strace \
        netcat-openbsd \
        procps \
        psmisc \
        time \
        less \
        nano \
        vim-tiny \
        man-db \
        manpages \
        manpages-dev \
        bash-completion \
    && rm -rf /var/lib/apt/lists/*

# Derslerin bir kopyası imajın içinde durur. Öğrenci boş bir klasör bağlarsa
# giriş scripti dersleri oraya kopyalar.
COPY . /opt/go-baslangic
RUN rm -rf /opt/go-baslangic/docker /opt/go-baslangic/.devcontainer

COPY --chmod=755 docker/ders.sh /usr/local/bin/ders
COPY --chmod=755 docker/giris.sh /usr/local/bin/giris

# Ölçüm dersleri dosyalarını /calisma'ya yazar. Burası bir Docker volume'üdür;
# Mac/Windows'taki yavaş paylaşımlı klasörden geçmez, gerçek Linux dosya
# sistemi ve page cache ölçülür.
# Linux'ta konteyner öğrencinin kendi kullanıcı kimliğiyle (--user) çalışır;
# bu yüzden ortak dizinler herkese yazılabilir olmalı.
RUN mkdir -p /calisma /kod /tmp/go-cache && chmod 1777 /calisma /kod /tmp/go-cache
ENV CALISMA=/calisma \
    GOCACHE=/tmp/go-cache
# Konteyner içinde Go'nun kendi sürüm indirmesini kapat.
ENV GOTOOLCHAIN=local

RUN printf '%s\n' \
    'PS1="\[\e[1;32m\]sysprog\[\e[0m\]:\[\e[1;34m\]\w\[\e[0m\]\$ "' \
    'alias ll="ls -lh"' \
    > /etc/profile.d/sysprog.sh \
    && echo '. /etc/profile.d/sysprog.sh' >> /etc/bash.bashrc

WORKDIR /kod
ENTRYPOINT ["giris"]
CMD ["bash"]

# ---------------------------------------------------------------------------
# Gerçek embedding modeli. Koordinatör bunu çocuk process olarak başlatır;
# pipe konteyner sınırını geçemeyeceği için model, koordinatörle aynı imajda olmalı.
FROM temel AS model

RUN apt-get update \
    && apt-get install -y --no-install-recommends python3 python3-venv \
    && rm -rf /var/lib/apt/lists/*

RUN python3 -m venv /opt/model-venv \
    && /opt/model-venv/bin/pip install --no-cache-dir \
        onnxruntime==1.31.0 tokenizers==0.23.3 numpy==2.5.4

# Model dosyaları belirli bir sürüme (commit) ve SHA256 özetine sabitlenir.
ARG E5_REV=614241f622f53c4eeff9890bdc4f31cfecc418b3
ADD --chmod=644 --checksum=sha256:ca456c06b3a9505ddfd9131408916dd79290368331e7d76bb621f1cba6bc8665 \
    https://huggingface.co/intfloat/multilingual-e5-small/resolve/${E5_REV}/onnx/model.onnx \
    /opt/modeller/e5-small/model.onnx
ADD --chmod=644 --checksum=sha256:0b44a9d7b51c3c62626640cda0e2c2f70fdacdc25bbbd68038369d14ebdf4c39 \
    https://huggingface.co/intfloat/multilingual-e5-small/resolve/${E5_REV}/tokenizer.json \
    /opt/modeller/e5-small/tokenizer.json

# İzinler ADD/COPY sırasında verilir. Sonradan "RUN chmod" yapmak dosyaları
# yeni bir katmana kopyalar ve modeller imajda iki kez yer kaplar.
COPY --chmod=755 docker/model/model-embedder.py /usr/local/bin/model-embedder
COPY --chmod=755 docker/model/model-kontrol.py /usr/local/bin/model-kontrol

# ---------------------------------------------------------------------------
# Cevap üreten dil modeli (bonus). llama.cpp resmî imajından alınır; içindeki
# CPU varyantları çalışma anında işlemciye göre (AVX2, AVX-512, ARM) seçilir.
FROM ghcr.io/ggml-org/llama.cpp:server-b11552 AS llamacpp

FROM model AS llm

RUN apt-get update \
    && apt-get install -y --no-install-recommends libgomp1 libssl3t64 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=llamacpp /app /opt/llama

ARG QWEN_REV=8fea620810c4afa23dd6443f999a48574c1611a3
ADD --chmod=644 --checksum=sha256:37ae482d336108d23516fa35e8e0c4126688d81018b87178a18d752a1357814f \
    https://huggingface.co/ggml-org/Qwen3.5-0.8B-GGUF/resolve/${QWEN_REV}/Qwen3.5-0.8B-Q8_0.gguf \
    /opt/modeller/qwen3.5-0.8b-q8_0.gguf

COPY --chmod=755 docker/model/model-uretici.sh /usr/local/bin/model-uretici
