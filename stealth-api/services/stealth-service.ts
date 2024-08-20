import App from "../app";

type SenderInfo = {
  privKey: string;
  pubKey: string;
};

type StealthInfo = {
  publicKey: string;
  address: string;
  viewTag: string;
};

type Info = {
  r: string;
};

type DiscoveredStealthInfo = {
  address: string;
  privKey: string;
};

interface Window {
    senderMeta: string;
    StealthPubKey: string;
    StealthAddress: string;
    StealthViewTag: string;
    DiscoveredStealthAddrs: string;
    DiscoveredStealthPrivKeys: string;
    recipientMeta: string;
}

declare const Go: any;

class StealthService {
    private app: App;

    private go = new Go();

    private window: Window;

    private _DB = {
        ids: ["Marija", "Mihajlo", "Mihailo", "Alex", "Malisa", "Milos"],
        idProcessed: -1,
        metaAddrRegistry: [] as Array<{
            id: string;
            R: string;
            stealthAddr: string;
        }>,

        txList: [{ senderPubKey: "", viewTag: "" }],
    };

    constructor(app: App, window: Window) {
        this.app = app;
        this.window = window;
    }

    public send(address: string, amount: number): void {
        if (amount < 0) {
            throw new Error("Invalid amount");
        }

        console.log(`Sending ${amount} to ${address}...`);
    }

    private async delay(ms: number): Promise<void> {
        return new Promise((res) => setTimeout(res, ms));
    }

    public async calculateStealth(
        setSenderInfo: (info: SenderInfo) => void,
        info: Info,
        setStealthInfo: (info: StealthInfo) => void,
        shortcircuit: boolean | null
    ): Promise<void> {
        if (shortcircuit == null) {
            WebAssembly.instantiateStreaming(
                fetch("ecpdksap.wasm"),
                this.go.importObject
            ).then((result) => {
                this.go.argv = ["js", "gen-send-info"];

                this.go.run(result.instance);

                const info2 = JSON.parse(this.window.senderMeta) as SenderInfo;

                console.log("generateSenderInfo", { info2 });

                setSenderInfo({ privKey: info2.privKey, pubKey: info2.pubKey });
            });
        }

        await this.delay(100);
        this.generateStealth(info, setStealthInfo);
    }

    private generateStealth(
        info: Info,
        setStealthInfo: (info: StealthInfo) => void
    ): void {
        WebAssembly.instantiateStreaming(
        fetch("ecpdksap.wasm"),
        this.go.importObject
        ).then((result) => {
            const info2 = JSON.parse(this.window.senderMeta) as { r: string };

            console.log("generateStealth", { info2 });

            info.r = info2.r;

            this.go.argv = ["js", "send", JSON.stringify(info)];
            this.go.run(result.instance);

            setStealthInfo({
                publicKey: this.window.StealthPubKey,
                address: this.window.StealthAddress,
                viewTag: this.window.StealthViewTag,
            });
        });
    }

    public receiveScan(
        info: { k: string; v: string },
        updateDiscoveredStealthInfo: (info: DiscoveredStealthInfo[]) => void
    ): void {
        WebAssembly.instantiateStreaming(
        fetch("ecpdksap.wasm"),
        this.go.importObject
        ).then((result) => {
            const Rs = this._DB.txList.map((el) => el.senderPubKey);
            const ViewTags = this._DB.txList.map((el) => el.viewTag);

            console.log({ info });

            this.go.argv = [
                "js",
                "receive-scan",
                JSON.stringify({
                k: info.k,
                v: info.v,
                Version: "v2",
                ViewTagVersion: "v0-1byte",
                Rs,
                ViewTags,
                }),
            ];
            this.go.run(result.instance);

            const dStealthAddrs = this.window.DiscoveredStealthAddrs.split(".");
            const dPrivKeys = this.window.DiscoveredStealthPrivKeys.split(".");

            const discovered = dStealthAddrs.map((address: any, k: any) => {
                return {
                    address,
                    privKey: dPrivKeys[k],
                };
            });

            updateDiscoveredStealthInfo(discovered);
        });
    }
    public async generateMetaRegistry(setLoaded: () => void): Promise<void> {
        WebAssembly.instantiateStreaming(
        fetch("ecpdksap.wasm"),
        this.go.importObject
        ).then((result) => {
            this.go.argv = ["js", "gen-recipient-info"];

            this.go.run(result.instance);

            const info = JSON.parse(this.window.recipientMeta) as {
                R: string;
                stealthAddr: string;
            };

            console.log("generateMetaRegistry", { info });

            this._DB.idProcessed += 1;

            this._DB.metaAddrRegistry.push({
                id: this._DB.ids[this._DB.idProcessed],
                ...info,
            });

            setLoaded();
        });
    }
}

export { StealthService };
