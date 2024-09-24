# Stealth API

This document provides an overview of the `API`, `BlockchainService` and `GoHandler` classes, which manage blockchain interactions and WebAssembly-based cryptographic operations. Also, `CLI` commands are supported for `API`.

## Table of Contents
- [Introduction](#introduction)
- [API Overview](#api-overview)
- [Blockchain Service](#blockchain-service)
- [Go Handler](#go-handler)
- [CLI Commands](#cli-commands)
- [Socket API](#socket-api)

## Introduction
This project is designed to interact with Ethereum-based smart contracts and handle cryptographic operations with WebAssembly. The main components include:

- **API**: Provide endpoints that allows users to interact with the blockchain and perdorm cryptographic operations via WebAssembly. 
- **Blockchain Service**: Manages interactions with the Ethereum blockchain, smart contracts, and event listeners.
- **Go Handler**: Interfaces with a WebAssembly (Wasm) module to perform cryptographic operations related to stealth addresses and meta-addresses.
- **CLI Commands**: Commands for running protocol functionalities through command line. 
- **Socket API**: This module defines WebSocket handlers for the Stealth API project. It uses the socket.io library to manage WebSocket connections and handle real-time communication between the server and clients.

---

## API Overview

The Stealth API is designed to provide an interface for interacting with stealth addresses in the context of the Ethereum blockchain. This `API` offers endpoints for generating and managing stealth addresses, retrieving metadata, and interacting with related cryptographic features. It sets up an Express server, configures middleware such as body-parser and query-parser, and exposes various routes using route handlers. This file is crucial for initializing the server and ensuring that incoming requests are processed correctly. Each route handler interacts with the database and services, ensuring appropriate responses and error handling for each request. The key routes and their functionalities:

1. **POST /register-address**
  - **Description**: Registers a new meta address
  - **Parameters**:
     - `id`: A unique identifier for the meta address.
     - `K`: The spending public key associated with the meta address.
     - `V`: The viewing public key associated with the meta address.
   - **Response**: Returns a confirmation message and the transaction hash.

3. **POST /send**
   - **Description**: Sends funds.
   - **Parameters**:
     - `recipientIdType`: Type of id
     - `id`: Optional parameter which represent recipient meta addres if recipientIdType is 'meta_address'
     - `recipientK`: Optional parameter which represent recipient spending key if recipientIdType isn't 'meta_address'.
     - `recipientV`: Optional parameter which represent recipient viewing key if recipientIdType isn't 'meta_address'.
     - `amount`: The amount of ETH to send.
     - `withProxy`: Bool value if you want to send with or without Proxy.
   - **Response**: Returns a confirmation message and the transaction receipt.

4. **GET /check-received**
   - **Description**: Retrieves a list of received transactions within specified block range.
   - **Parameters**:
     - `fromBlock`: Optional parameter which represent first block service need to check for transaction.
     - `toBlock`: Optional parameter which represent last block service need to check for transaction.
   - **Response**: Returns an array of received transactions with details.

5. **POST /transfer/:receiptId**
   - **Description**: Transfers funds from a received transaction based on the provided receipt ID.
   - **Parameters**:
     - `receiptId`: The ID of the transaction receipt.
     - `address`: The address where you want to transfer funds
     - `amount`: The amount of funds you want to transfer
   - **Response**: Returns a confirmation message and the transaction receipt.

## Blockchain Service

This service acts as a bridge between the application and the blockchain, ensuring smooth interactions and logging relevant events. The main
features of `BlockchainService` class includes:

1. **registerMetaAddress(id: string, K: string, V: string)**
   - **Description**: Registers a new meta address on the blockchain.
   - **Parameters**:
     - `id`: Unique identifier for the meta address.
     - `K`: Spending public key associated with the meta address.
     - `V`: Viewing public key associated with the meta address.
   - **Returns**: The transaction receipt after registration.
   - **Error Handling**: Logs errors encountered during registration.

2. **resolveMetaAddress(id: string)**
   - **Description**: Resolves a meta address to retrieve its associated K and V values.
   - **Parameters**:
     - `id`: The ID of the meta address to resolve.
   - **Returns**: An object containing K and V values.
   - **Error Handling**: Logs errors if resolution fails.

3. **sendEthViaProxy(stealthAddress: string, R: string, viewTag: string, amount: string)**
   - **Description**: Sends ETH to a stealth address via a proxy.
   - **Parameters**:
     - `stealthAddress`: The stealth address to send ETH to.
     - `R`: The ephemeral public key.
     - `viewTag`: The view tag.
     - `amount`: The amount of ETH to send.
   - **Returns**: The transaction receipt.
   - **Error Handling**: Logs errors during the sending process.

4. **ethSentWithoutProxy(stealthAddress: string, R: string, viewTag: string, amount: string)**
   - **Description**: Sends ETH directly to a stealth address without using a proxy.
   - **Parameters**: Same as `sendEthViaProxy`.
   - **Returns**: The transaction receipt.
   - **Error Handling**: Logs errors during the sending process.

5. **listenMetaAddressRegistredEvent()**
   - **Description**: Listens for the `MetaAddressRegistered` event emitted by the blockchain and logs the details.
   - **Returns**: None.
   - **Logs**: Logs when the event listener is active.

5. **listenAnnouncementEvent()**
   - **Description**: Listens for the `Announcement` event and processes relevant data, saving to the database if conditions are met.
   - **Returns**: None.
   - **Logs**: Logs when the event listener is active.

7. **transferEth(address: string, amount: string, privKey: string)**
    - **Description**: Transfers ETH from the wallet using the provided private key.
    - **Parameters**:
      - `address`: The recipient's address.
      - `amount`: The amount of ETH to transfer.
      - `privKey`: The private key of the sender's wallet.
    - **Returns**: The transaction receipt.
    - **Error Handling**: Logs errors if the transfer fails.

## Go Handler

The `GoHandler` class is essential for integrating Go-based cryptography operations into the JavaScript environment, leveraging WebAssembly for performance and security. The main features of this class includes:

1. **genSenderInfo()**
   - **Description**: Generates sender information using the WebAssembly module.
   - **Returns**: A promise that resolves to an `Info` object containing sender details.
   - **Error Handling**: Rejects the promise if instantiation or execution fails.

2. **genRecipientInfo()**
   - **Description**: Generates recipient information using the WebAssembly module.
   - **Returns**: A promise that resolves to an `Info` object containing recipient details.
   - **Error Handling**: Rejects the promise if instantiation or execution fails.

3. **send(r: string, K: string, V: string)**
   - **Description**: Sends information using the WebAssembly module and returns the stealth address and view tag.
   - **Parameters**:
     - `r`: The ephemeral key.
     - `K`: The spending public key for sending.
     - `V`: The viewing public key for sending.
   - **Returns**: A promise that resolves to a `SendInfo` object with stealth address and view tag.
   - **Error Handling**: Rejects the promise if instantiation or execution fails.

4. **receiveScan(k: string, v: string, Rs: string[], viewTags: string[])**
   - **Description**: Scans for received funds using the WebAssembly module.
   - **Parameters**:
     - `k`: The spending private key for scanning.
     - `v`: The viewing private key for scanning.
     - `Rs`: An array of ephemeral keys to scan.
     - `viewTags`: An array of view tags.
   - **Returns**: A promise that resolves to an array of `ReceiveScanInfo` objects containing discovered addresses and private keys.
   - **Error Handling**: Rejects the promise if instantiation or execution fails.

## CLI Commands

The `CLI` is designed to facilitate interaction with the API directly from the command line, making it easier to perform common operations. Below
are the available commands along with their options:

#### 1. `register-address`
- **Description**: Register a new meta address on the contract.
- **Options**:
  - `--id <string>`: The ID associated with the meta address.
  - `--K <string>`: The spending public key for the meta address.
  - `--V <string>`: The viewing public key for the meta address.
- **Usage**: 
  ```bash
    npm run cli -- register-address --id "id_value" --K "K_value" --V "V_value"
  ```

#### 2. `send`
- **Description**: Sends funds to a specified recipient's stealth address.
- **Options**:
  - `recipientIdType <string>`: Specifies the type of recipient ID, e.g., meta_address.
  - `id <string>`: Optional ID
  - `recipientK <string>`: Optional spending public key for the recipient (if recipientIdType is meta_address).
  - `recipientV <string>`: Optional viewing public key for the recipient (if recipientIdType is meta_address).
  - `amount <string>`: The amount of funds to send.
  - `withProxy <bool>`: An optional flag to indicate whether to send via a proxy.
- **Usage**: 
  ```bash
    npm run cli -- send --recipientIdType "idType_value" --id "id_value" --recipientK "K_value" --recipientV "V_value" --amount "amount_value" --withProxy
  ```

#### 3. `check-received`
- **Description**: Checks for any received transactions between specified block ranges.
- **Options**:
  - `fromBlock <number>`: The block number to start the query from (default: 0).
  - `toBlock <number>`: The block number to end the query at (default: the latest block).
- **Usage**: 
  ```bash
    npm run cli -- check-received --fromBlock fromBlock --toBlock toBlock
  ```

#### 4. `transfer`
- **Description**: Transfer received funds to another address.
- **Options**:
  - `receiptId <number>`: The ID of the receipt for which the transfer will occur.
  - `address <string>`: Optional transfer address. If not provided, it will use the default address from the configuration.
  - `amount <number>`: Optional transfer amount. If not provided, the default amount is set to 0.001.
- **Usage**: 
  ```bash
    npm run cli -- transfer --receiptId receiptId --address "address" --amount amount
  ```

## Socket API

This module defines WebSocket handlers for the Stealth API project. It uses the socket.io library to manage WebSocket connections and handle real-time communication between the server and clients. The module allows clients to interact with the blockchain. The main function, `setupSocketHandlers`, listens for specific WebSocket events and invokes the appropriate business logic in the application.

The `setupSocketHandlers` function sets up various WebSocket events for handling client requests. It takes two arguments:
- `io: Server` - The socket.io server instance.
- `app: App` - The main application instance, which contains the business logic and services for blockchain, database interactions, and logging.
The main ecents of this functions are:

#### 1. register-address
- **Description**: Register a new meta address on the contract.
- **Parameters**:
  - `data`:
    - `id`: A unique identifier for the meta address.
    - `K`: The spending public key associated with the meta address.
    - `V`: The viewing public key associated with the meta address.
  - `callback`:
    - `message`
    - `id` 

#### 2. send
- **Description**: Sends funds to a specified recipient's stealth address.
- **Parameters**:
  - `data`:
    - `recipientIdType`: Type of id
    - `id`: Optional parameter which represent recipient meta addres if recipientIdType is 'meta_address'
    - `recipientK`: Optional parameter which represent recipient spending key if recipientIdType isn't 'meta_address'.
    - `recipientV`: Optional parameter which represent recipient viewing key if recipientIdType isn't 'meta_address'.
    - `amount`: The amount of ETH to send.
    - `withProxy`: Bool value if you want to send with or without Proxy.
  - `callback`:
    - `message`
    - `data` = { `stealthAddress`, `ephemeralPubKey`, `viewTag`, `amount` }

#### 3. check-received
- **Description**: Checks for any received transactions between specified block ranges.
- **Parameters**:
  - `data`:
    - `fromBlock`: Optional parameter which represent first block service need to check for transaction.
    - `toBlock`: Optional parameter which represent last block service need to check for transaction.
  - `callback`:
    - `message`
    - `receipt`

#### 4. transfer
- **Description**: Transfer received funds to another address.
- **Parameters**:
  - `data`:
    - `receiptId`: The ID of the transaction receipt.
    - `address`: The address where you want to transfer funds
    - `amount`: The amount of funds you want to transfer
  - `callback`:
    - `message`
    - `transaction`