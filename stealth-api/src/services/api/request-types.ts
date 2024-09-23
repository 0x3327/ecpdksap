export type SendFundsRequest = {
    recipientIdType: 'eth_ens' | 'address' | 'meta_address',
    ens?: string,
    address?: string,
    recipientK?: string,
    recipientV?: string,
    amount: number,
    withProxy: boolean
}

export type CheckReceivedRequest = {
    fromBlock: number,
    toBlock: number,
}

export type TransferReceivedFundsRequest = {
    address: string,
    amount: number,
}

export type RegisterUser = {
    pid: number,
    name: string,
    privateKey: string,
    publicKeyX: string,
    publicKeyY: string,
}