export type Info = {
    k: string;
    v: string;
    r: string;
    K: string;
    V: string;
    R: string;
    P_Sender: string;
    viewTag: string;
    P_Recipient: string;
    Version: string;
    ViewTagVersion: string;
}

export type SendInfo = {
    pubKey: string;
    address: string;
    viewTag: string;
}

export type ReceiveScanInfo = {
    address: string;
    privKey: string;
}