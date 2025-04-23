## ECPDKSAP Contracts: Overview

There exist two main contracts:

- `ECPDKSAP_MetaAddressRegistry`:

  - used to register recipient's meta addresses (raw bytes) using a human-friendly ID (i.e. `vitalik.eth`)
  - holds the data needed to resolve the meta address from the ID

- `ECPDKSAP_Announcer`:

  - used to notify the recipient(s) of new potential eth transfers to the stealth address they control
  - note: calls the singleton `ERC5564Announcer` contract (see: [EIP-5564](https://eips.ethereum.org/EIPS/eip-5564))

### Detailed Flow:

1. **Recipient** calls `ECPDKSAP_MetaAddressRegistry.registerMetaAddress(string memory _id, bytes memory _metaAddress)`

- where `_id` is a human-readable text that corresponds to the underlying `_metaAddress`

2. (Optional) **Sender** calls `ECPDKSAP_MetaAddressRegistry.resolve(string memory _id)`

- resolves the `_id` to `bytes memory _metaAddress`

3. **Sender** generate Recipient's stealth address (Ethereum EOA) using the bytes `_memoryAddress` and the typescript API lib.

4. **Sender** calls `ECPDKSAP_Announcer` with either:

- `.sendEthViaProxy(address payable _stealthAddress, bytes memory _R, bytes memory _viewTag)`
  - Sends funds to `_stealthAddress` using the announcer contract as a proxy
  - Notifies the recipient(s) about the transfers using event emission
- `.ethSentWithoutProxy(bytes memory _R, bytes memory _viewTag)`
  - Notifies the recipient(s) about the transfers using event emission
  - _Note: the actual ETH transfer to the generated `_stealthAddress` can happen from a different account (not the `msg.sender` for this contract call )_

5. **Recipient** catches all the events emitted and parses through them, potentially generating their stealth addresses and corresponding private keys

## Developer: Getting Started

Prepare the dev. environment (& substitute variable values):

```
cp .env.example .env
```

_Note: Due to calling the `ERC5564Announcer` singleton contract, **testing** is done using a Sepolia fork._

```
source .env && forge test --fork-url $SEPOLIA_RPC_URL
```

Deployment:

```
forge script --chain sepolia script/Deployment.s.sol:DeploymentScript --rpc-url $SEPOLIA_RPC_URL --broadcast --verify -vvvv
```

## Deployed Contracts (Sepolia)

[ECPDKSAP_MetaAddressRegistry](https://sepolia.etherscan.io/address/0xb4b82918613524db74967ca6c71979cd030b7991)

- ```
  0xB4B82918613524DB74967CA6c71979cD030B7991
  ```

[ECPDKSAP_Announcer](https://sepolia.etherscan.io/address/0x79820c9a124023d47bbcc6d0a24db4d0075ca724)

- ```
  0x79820C9a124023D47BbCC6d0a24DB4D0075Ca724
  ```

Existing singleton [ERC5564Announcer](https://sepolia.etherscan.io/address/0x55649E01B5Df198D18D95b5cc5051630cfD45564)

- ```
  0x55649E01B5Df198D18D95b5cc5051630cfD45564
  ```
