#!/opt/model-venv/bin/python3
"""Model kurulumunu doğrular: embedder protokolü ve (varsa) üretici.

Kullanım: model-kontrol
"""
import json
import os
import shutil
import struct
import subprocess
import sys
import time
import urllib.request


# Bonus (ANSWER) için önerilen sistem mesajı. Küçük modeller bağlam dışına
# taşmaya meyillidir; kısıtı açıkça yazmak cevapları belirgin şekilde iyileştirir.
SISTEM = ("Sen bir soru cevaplama asistanısın. Yalnızca verilen BAĞLAM'daki bilgileri kullan. "
          "Cevap bağlamda yoksa sadece 'Bu bilgi depoda yok.' yaz. Kısa ve düzgün Türkçe ile cevap ver.")


def oku_tam(akis, n):
    """Pipe'tan tam n bayt oku. Tek read() n bayttan az döndürebilir!"""
    parca = b""
    while len(parca) < n:
        ek = akis.read(n - len(parca))
        if not ek:
            raise EOFError(f"{len(parca)}/{n} bayt okundu, akış kapandı")
        parca += ek
    return parca


def embedder_kontrol():
    print("== Embedder (multilingual-e5-small)")
    basla = time.monotonic()
    p = subprocess.Popen(["model-embedder"], stdin=subprocess.PIPE, stdout=subprocess.PIPE)
    hazir = p.stdout.readline().decode().strip()
    print(f"   açılış: {hazir!r} ({time.monotonic() - basla:.1f} sn)")
    assert hazir == "HAZIR 384", "beklenen: HAZIR 384"

    cumleler = [
        "İstanbul Boğazı Avrupa ile Asya'yı ayırır.",
        "iki kıtayı birbirinden ayıran su yolu",
        "Menemen domates, biber ve yumurtayla yapılır.",
    ]
    vektorler = []
    for c in cumleler:
        t = time.monotonic()
        p.stdin.write(c.encode() + b"\n")
        p.stdin.flush()
        ham = oku_tam(p.stdout, 1536)
        v = struct.unpack("<384f", ham)
        vektorler.append(v)
        norm = sum(x * x for x in v) ** 0.5
        print(f"   {c[:45]:45s} -> 1536 bayt, norm={norm:.4f}, {1000 * (time.monotonic() - t):.0f} ms")

    def ic(a, b):
        return sum(x * y for x, y in zip(a, b))

    yakin, uzak = ic(vektorler[0], vektorler[1]), ic(vektorler[0], vektorler[2])
    print(f"   benzerlik: boğaz~kıtalar={yakin:.3f}  boğaz~menemen={uzak:.3f}")
    assert yakin > uzak, "anlamca yakın cümleler daha benzer çıkmalıydı"

    p.stdin.close()
    p.wait(timeout=10)
    print(f"   stdin kapatıldı, embedder çıkış kodu: {p.returncode}")


def uretici_kontrol():
    if not shutil.which("model-uretici"):
        print("== Üretici: bu imajda yok (sadece :llm imajında). Atlandı.")
        return
    print("== Üretici (Qwen3.5-0.8B, llama.cpp)")
    port = 18081
    p = subprocess.Popen(["model-uretici", str(port)], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    try:
        basla = time.monotonic()
        while True:
            try:
                urllib.request.urlopen(f"http://127.0.0.1:{port}/health", timeout=1)
                break
            except Exception:
                if p.poll() is not None or time.monotonic() - basla > 120:
                    raise RuntimeError("üretici açılmadı")
                time.sleep(0.5)
        print(f"   açılış: {time.monotonic() - basla:.1f} sn")
        istek = {
            "messages": [
                {"role": "system", "content": SISTEM},
                {"role": "user", "content": "BAĞLAM:\n- fsync sistem çağrısı dosyadaki değişikliklerin "
                                            "diske kalıcı olarak yazılmasını sağlar.\n\n"
                                            "SORU: Elektrik kesilince veri kaybolmasın diye ne yapmalıyım?"},
            ],
            "temperature": 0, "max_tokens": 80,
        }
        t = time.monotonic()
        req = urllib.request.Request(f"http://127.0.0.1:{port}/v1/chat/completions",
                                     json.dumps(istek).encode(), {"Content-Type": "application/json"})
        cevap = json.load(urllib.request.urlopen(req, timeout=300))
        tm = cevap.get("timings", {})
        print(f"   cevap: {cevap['choices'][0]['message']['content'].strip()}")
        print(f"   {time.monotonic() - t:.1f} sn, üretim {tm.get('predicted_per_second', 0):.1f} token/sn")
    finally:
        p.terminate()
        p.wait(timeout=10)


if __name__ == "__main__":
    print(f"İşlemci çekirdeği: {os.cpu_count()}")
    embedder_kontrol()
    uretici_kontrol()
    print("TAMAM")
