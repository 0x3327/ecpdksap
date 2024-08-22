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
                readFileSync(path.join(__dirname, '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "gen-send-info"];

                this.go.run(result.instance);

                // console.log("global genSenderInfo", global);
                const info = JSON.parse((global as any).senderMeta);
                resolve(info);
            }).catch(err => reject(err))}) 
    };

    genRecipientInfo() {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "gen-recipient-info"];

                this.go.run(result.instance);

                // console.log("global genSenderInfo", global);
                const info = JSON.parse((global as any).recipientMeta);

                this._DB.idProccessed += 1;

                this._DB.metaAddrRegistry.push({ id: this._DB.ids[this._DB.idProccessed], ...info});

                resolve(info);
            }).catch(err => reject(err))})
    };

    send(recipientInfo: Info) {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                const senderInfo = JSON.parse((globalThis as any).senderMeta);
                // console.log("info", info);
                recipientInfo.r = senderInfo.r;
                // console.log("senderInfo", senderInfo);
                // console.log("info", info);
                
                this.go.argv = ["js", "send", JSON.stringify(recipientInfo)];

                this.go.run(result.instance);

                const stealthPubKey = (global as any).StealthPubKey;
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

    receiveScan(recipientInfo: Info, sendInfo: SendInfo) {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                // const Rs = this._DB.txList.map((el) => el.senderPubKey);
                // const ViewTags = this._DB.txList.map((el) => el.viewTag);

                // console.log("global", global);
                console.log("global StealthPubKey", (globalThis as any).StealthPubKey);
                console.log("global StealthViewTag", (globalThis as any).StealthViewTag);
                const Rs = [(globalThis as any).StealthPubKey];
                // const ViewTags = [(globalThis as any).StealthViewTag];

                this.go.argv = ["js", "receive-scan", JSON.stringify({
                    k: recipientInfo.k,
                    v: recipientInfo.v,
                    Version: "v2",
                    ViewtagVersion: "v0-1byte",
                    Rs,
                    ViewTags: ["2c"],
                })];

                this.go.run(result.instance);

                console.log("prosao run");

                console.log("discoveredSA", (global as any).DiscoveredStealthAddrs);
                console.log("discoveredPK", (global as any).DiscoveredStealthPrivKeys);
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