## Overview

The repository contains implementation code for all [ECPDKSAP protocol versions](./docs).

All protocol variations are implemented in Go programming language and use [consensys/gnark-crypto](https://github.com/Consensys/gnark-crypto) library.

<<<<<<< HEAD
=======
Forked version: [0x3327/gnark-crypto](https://github.com/0x3327/gnark-crypto) containing custom function implementations needed by ECPDKSAP protocols.

>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
Smart contracts are developed using Foundry framework and follow the [EIP-5564](https://eips.ethereum.org/EIPS/eip-5564) standard.

## Project structure

The project's dir. structure is as following:

- `./impl`: Off-chain protocol code (used by sender and recipient to generate/scan data)
- `./sc`: On-chain contracts (used for sender - recipient "communication")
- `./backend`: Event indexers, API modules, ...

## Additional resources

Original paper: https://arxiv.org/abs/2312.12131
