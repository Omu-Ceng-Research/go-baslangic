#!/opt/model-venv/bin/python3
"""Gerçek embedding modeli: multilingual-e5-small (384 boyut, Türkçe destekli).

Ödevdeki embedder pipe protokolünü konuşur:
  1. Açılışta stdout'a tek satır:          HAZIR 384\n
  2. stdin'den bir satır metin okur:        <metin>\n
  3. stdout'a TAM 1536 bayt yazar:          384 x float32, little-endian, birim uzunlukta
  4. stdin kapanınca (EOF) çıkar.

stdout yalnızca protokol içindir. Tüm loglar stderr'e gider.

Ortam değişkenleri:
  EMBEDDER_THREADS   ONNX Runtime iş parçacığı sayısı (varsayılan 2)
  EMBEDDER_MODEL     model klasörü (varsayılan /opt/modeller/e5-small)
"""
import os
import sys
import time

BOYUT = 384
ONEK = "query: "  # e5 modelleri önek bekler; kayıt ve soru için aynısı kullanılır
MAKS_TOKEN = 512


def log(*args):
    print("[model-embedder]", *args, file=sys.stderr, flush=True)


def main():
    basla = time.monotonic()
    # ONNX Runtime HOME'a yazmaya çalışır; root olmayan kullanıcıda HOME=/ olabilir.
    if not os.access(os.environ.get("HOME", "/"), os.W_OK):
        os.environ["HOME"] = "/tmp"
    import numpy as np
    import onnxruntime as ort
    from tokenizers import Tokenizer

    ort.set_default_logger_severity(3)  # yalnızca hatalar
    ort.disable_telemetry_events()

    klasor = os.environ.get("EMBEDDER_MODEL", "/opt/modeller/e5-small")
    threads = int(os.environ.get("EMBEDDER_THREADS", "2"))

    tok = Tokenizer.from_file(os.path.join(klasor, "tokenizer.json"))
    tok.enable_truncation(MAKS_TOKEN)

    opts = ort.SessionOptions()
    opts.intra_op_num_threads = threads
    opts.inter_op_num_threads = 1
    oturum = ort.InferenceSession(os.path.join(klasor, "model.onnx"), opts,
                                  providers=["CPUExecutionProvider"])

    def embed(metin):
        enc = tok.encode(ONEK + metin)
        ids = np.array([enc.ids], dtype=np.int64)
        maske = np.array([enc.attention_mask], dtype=np.int64)
        cikti = oturum.run(None, {"input_ids": ids, "attention_mask": maske,
                                  "token_type_ids": np.zeros_like(ids)})[0][0]
        m = maske[0][:, None].astype(np.float32)
        v = (cikti * m).sum(0) / m.sum()  # ortalama havuzlama (mean pooling)
        v /= np.linalg.norm(v)
        return v.astype("<f4").tobytes()  # little-endian float32

    embed("ısınma")  # ilk çağrı yavaştır; HAZIR demeden önce yapalım
    log(f"model yüklendi ({time.monotonic() - basla:.1f} sn, {threads} thread, pid {os.getpid()})")

    cikis = sys.stdout.buffer
    cikis.write(f"HAZIR {BOYUT}\n".encode())
    cikis.flush()

    for satir in sys.stdin.buffer:  # EOF'ta döngü biter
        metin = satir.rstrip(b"\r\n").decode("utf-8", errors="replace")
        vektor = embed(metin)
        assert len(vektor) == BOYUT * 4
        cikis.write(vektor)
        cikis.flush()

    log("stdin kapandı, çıkıyorum")


if __name__ == "__main__":
    try:
        main()
    except (KeyboardInterrupt, BrokenPipeError):
        pass
