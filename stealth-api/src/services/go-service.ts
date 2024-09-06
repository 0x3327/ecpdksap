import { readFileSync } from 'fs';
import path from 'path';
import { Info, ReceiveScanInfo, SendInfo } from '../types';
import App from '../app';

class GoHandler {
    go: any;

    constructor(app: App) {
        const { Go } = (globalThis as any);
        this.go = new Go()
    }

    genSenderInfo(): Promise<Info> {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "gen-send-info"];

                this.go.run(result.instance);

                const info = JSON.parse((global as any).senderMeta);
                resolve(info);
            }).catch(err => reject(err))})
    };

    genRecipientInfo(): Promise<Info> {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "gen-recipient-info"];

                this.go.run(result.instance);

                const info = JSON.parse((global as any).recipientMeta);

                resolve(info);
            }).catch(err => reject(err))})
    };

    send(r: string, K: string, V: string): Promise<SendInfo> {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "send", JSON.stringify({
                    r: r,
                    K: K,
                    V: V,
                    Version: "v2",
                    ViewTagVersion: "v0-1byte",
                })];

                this.go.run(result.instance);

                const stealthPubKeyCleaned = (global as any).StealthPubKey.replace("E([", "").replace("])", "");
                let [x, y] = stealthPubKeyCleaned.split(",");
                const stealthPubKey = `${x}.${y}`;
                (global as any).StealthPubKey = stealthPubKey;
                const stealthAddress = (global as any).StealthAddress;
                const stealthViewTag = (global as any).StealthViewTag;

                const sendInfo: SendInfo = {
                    pubKey: stealthPubKey,
                    address: stealthAddress,
                    viewTag: stealthViewTag
                };
                resolve(sendInfo);
            })
        })
    };

    receiveScan(k: string, v: string, Rs: string[], viewTags: string[]): Promise<ReceiveScanInfo[]> {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "receive-scan", JSON.stringify({
                    k: k,
                    v: v,
                    Rs,
                    Version: "v2",
                    ViewTags: viewTags,
                    ViewtagVersion: "v0-1byte",
                })];

                this.go.run(result.instance);

                const dStealthAddrs = (global as any).DiscoveredStealthAddrs.split(".");
                const dPrivKeys = (global as any).DiscoveredStealthPrivKeys.split(".");
                const receiveScanInfo: ReceiveScanInfo[] = dStealthAddrs.map((address: any, k: any) => {
                    return {
                        address,
                        privKey: dPrivKeys[k],
                    };
                });

                resolve(receiveScanInfo);
            })
        })
    };

}

export default GoHandler;