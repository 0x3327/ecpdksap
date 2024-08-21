import { readFileSync } from 'fs';
import path from 'path';

console.log(global);

class GoHandler {
    go: any;

    constructor() {

        const { Go } = (globalThis as any);
        this.go = new Go()
    }

    genSenderInfo() {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "gen-send-info"];
                this.go.run(result.instance);
                const info2 = JSON.parse((global as any).senderMeta);
                resolve(info2)
            }).catch(err => reject(err))}) 
    }

    genRecipientInfo() {
        return new Promise((resolve, reject) => {
            WebAssembly.instantiate(
                readFileSync(path.join(__dirname, '..', 'public', 'ecpdksap.wasm')),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "gen-recipient-info"];

                this.go.run(result.instance);

                const info = JSON.parse((globalThis as any).recipientMeta);
                resolve(info)
            }).catch(err => reject(err))})
    }

}

export default GoHandler;