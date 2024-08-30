class Profile {
    private k: string;
    private v: string;

    constructor(k: string, v: string) {
        this.k = k;
        this.v = v;
    }
    
    getK(): string {
        return this.k;
    }

    getV(): string {
        return this.v;
    }
}

export default Profile;
