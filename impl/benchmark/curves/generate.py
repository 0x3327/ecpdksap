import os

CURVES = [
    {"name": "bw6-761"},
    {"name": "bls12-377"},
    {"name": "bls12-381"},
    {"name": "bls24-315"},
    {"name": "bw6-633"},
]

for c in CURVES:

    text = open('./bn254/main.go').read()

    os.makedirs(c['name']) if not os.path.exists(c['name']) else None

    with open(f"{c['name']}/main.go", "w") as outFile:
        
        text = text.replace("gnark-crypto/ecc/bn254", f"gnark-crypto/ecc/{c['name']}")
        text = text.replace("Running `bn254` Benchmark", f"Running `{c['name']}` Benchmark ")
        
        if c['name'] == "bls24-315":
            text = text.replace("S.C0.B0.A0.BigInt(b_asBigInt)", "S.D0.C0.B0.A0.BigInt(b_asBigInt)")
        
        elif c['name'] == 'bw6-633' or c['name'] == 'bw6-761':
            text = text.replace("S.C0.B0.A0.BigInt(b_asBigInt)", "S.B0.A0.BigInt(b_asBigInt)")

        outFile.write(text)