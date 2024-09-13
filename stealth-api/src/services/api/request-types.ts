export type SendFundsRequest = {
    recipientIdType: 'id' | 'eth_ens' | 'address' | 'meta_address',
    id?: string,
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