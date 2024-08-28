## ECPDKSAP Contracts: Overview

There exit two main contracts:

- `ECPDKSAP_MetaAddressRegistry`:

  - used to register recipient's meta addresses (raw bytes) using a people-friendly ID (i.e. `vitalik.eth`)
  - holds the data needed to resolve the meta address from the ID

- `ECPDKSAP_Announcer`:
  - used to notify the recipient(s) of new potential eth transfers to the stealth address they control
  - note: calls the singleton `ERC5564Announcer` contract (see: [EIP-5564](https://eips.ethereum.org/EIPS/eip-5564))

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
TODO
```
