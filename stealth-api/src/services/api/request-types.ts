export type SendFundsRequest = {
    recipientIdType: 'eth_dns' | 'address',
    dns?: string,
    address?: string,
    senderr?: string,
    recipientK?: string,
    recipientV?: string,
    amount: number,
}

export type CheckReceivedRequest = {
    fromBlock: number,
    toBlock: number,
}

export type TransferReceivedFundsRequest = {
    address: string,
    amount: number,
}