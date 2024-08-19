import App from '../app';

class StealthService {
    private app: App;

    constructor(app: App) {
        this.app = app;
    }

    public send(address: string, amount: number): void {
        if (amount < 0) {
            throw new Error('Invalid amount');
        }

        console.log(`Sending ${amount} to ${address}...`);
    }
}

export { StealthService };
