import { readFileSync } from 'fs';
import path from 'path';

type Info = {
    k: string;
    v: string;
    r: string;
    K: string;
    V: string;
    R: string;
    P_sender: string;
    viewTag: string;
    P_Recipient: string;
    Version: string;
    ViewTagVersion: string;
};

type SendInfo = {
    pubKey: string;
    address: string;
    viewTag: string;
}

type MetaAddr = {
    id: string;
    info: Info;
};

type txListType = {
    senderPubKey: string;
    viewTag: string;
}

type DB = {
    ids: string[];
    idProccessed: number;
    metaAddrRegistry: MetaAddr[];
    txList: txListType[];
}

class GoHandler {
    go: any;

    constructor() {
        const { Go } = (globalThis as any);
        this.go = new Go()
    }

    private _DB: DB = {
        ids: ["Marija", "Mihajlo", "Mihailo", "Alex", "Malisa", "Milos"],
        idProccessed: -1,
        metaAddrRegistry: [],
      
        txList: [{ senderPubKey: "", viewTag: "" }],
    };

    genSenderInfo() {
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

    genRecipientInfo() {
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

    send(): Promise<SendInfo> {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                const recipientInfo = JSON.parse((global as any).recipientMeta);
                const senderInfo = JSON.parse((global as any).senderMeta);
                
                this.go.argv = ["js", "send", JSON.stringify({
                    r: senderInfo.r,
                    K: recipientInfo.K,
                    V: recipientInfo.V,
                    Version: recipientInfo.Version,
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

    receiveScan() {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                const recipientInfo = JSON.parse((global as any).recipientMeta);
                const senderInfo = JSON.parse((global as any).senderMeta);

                const Rs = [senderInfo.R];
                const ViewTags = [(global as any).StealthViewTag];

                this.go.argv = ["js", "receive-scan", JSON.stringify({
                    k: recipientInfo.k,
                    v: recipientInfo.v,
                    Rs,
                    Version: "v2",
                    ViewTags,
                    ViewtagVersion: "v0-1byte",
                })];

                this.go.run(result.instance);

                const dStealthAddrs = (global as any).DiscoveredStealthAddrs.split(".");
                const dPrivKeys = (global as any).DiscoveredStealthPrivKeys.split(".");
                const receiveScanInfo = dStealthAddrs.map((address: any, k: any) => {
                    return {
                        address,
                        privKey: dPrivKeys[k],
                    };
                });
                resolve(receiveScanInfo);
            })
        })
    };

    // genExample() {
    //     return new Promise((resolve, reject) => {
    //         WebAssembly.instantiate(
    //             readFileSync(path.join(__dirname, '..', 'public', 'ecpdksap.wasm')),
    //             this.go.importObject
    //         ).then((result) => {
    //             this.go.argv = ["js", "gen-example"];

    //             this.go.run(result.instance);
    //             const info = JSON.parse((globalThis as any).Example)
    //             resolve(info);
    //         })
    //     })
    // }

}

export default GoHandler;